package vkscraper

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/kandayo/vkscraper/pkg/vk"
)

// TODO: check if vk sends a Retry-After header
const defaultRateLimitSleepTime = 3 * time.Second

// TODO: make it customizable
const timeout = time.Minute

type Instance struct {
	Config    Config
	Vk        *vk.Client
	Log       *log.Logger
	transport *http.Client
}

type Config struct {
	BaseDir    string
	FastUpdate bool
}

type Context struct {
	ScreenName       string
	UserID           int
	DirectoryHandler func(values ...interface{}) string
}

func New(client *vk.Client, config Config) *Instance {
	log := log.New(os.Stdout, "", 0)
	return &Instance{
		Config: config,
		Vk:     client,
		Log:    log,
		transport: &http.Client{
			Timeout: timeout,
		},
	}
}

func (instance *Instance) DownloadProfiles(screenNames []string) {
	for _, screenName := range screenNames {
		err := instance.DownloadProfile(screenName)
		if err != nil {
			instance.Log.Printf("[BUG] Could not download profile %s, got error: %s\n", screenName, err)
			continue
		}
	}
}

func (instance Instance) DownloadProfile(screenName string) error {
	userInfo, err := instance.Vk.ResolveScreenName(screenName)
	if err != nil {
		ok := handleApiError(err)
		if !ok {
			return err
		}
	}

	userID := userInfo.ID
	if userInfo.Type == "group" {
		// Negative IDs are used to designate a community ID.
		userID = userID * -1
	}

	profileDirectoryHandler := func(values ...interface{}) string {
		parsedValues := []string{instance.Config.BaseDir, screenName}
		for _, value := range values {
			switch t := value.(type) {
			case string:
				parsedValues = append(parsedValues, t)
			case int:
				parsedValues = append(parsedValues, strconv.Itoa(t))
			}
		}
		return filepath.Join(parsedValues...)
	}

	// TODO: Refactor
	ctx := Context{
		ScreenName:       screenName,
		UserID:           userID,
		DirectoryHandler: profileDirectoryHandler,
	}

	instance.Log.Printf("Downloading profile %s (vkid: %d)\n", screenName, userID)

	_ = createDirectory(ctx.DirectoryHandler("stories"))
	err = instance.DownloadStories(ctx)
	if err != nil {
		instance.Log.Printf("Error while retrieving stories from %s (vkid %d): %s\n", ctx.ScreenName, ctx.UserID, err)
	}

	_ = createDirectory(ctx.DirectoryHandler("photos"))
	err = instance.downloadPhotos(ctx)
	if err != nil {
		instance.Log.Printf("Error while retrieving posts from %s (vkid %d): %s\n", ctx.ScreenName, ctx.UserID, err)
	}

	return nil
}

func (instance *Instance) DownloadStories(ctx Context) error {
	collection, err := instance.Vk.StoriesGet(ctx.UserID)
	if err != nil {
		ok := handleApiError(err)
		if !ok {
			return err
		}
	}
	for _, story := range collection {
		downloaded, err := instance.downloadStoryItem(ctx, story)
		if err != nil {
			return err
		}
		if !downloaded && instance.Config.FastUpdate {
			return nil
		}
	}
	return nil
}

// TODO: duplicated code
func (instance *Instance) downloadStoryItem(ctx Context, story vk.Story) (bool, error) {
	if story.IsPhoto() {
		return instance.downloadStoryPhoto(ctx, story)
	} else {
		return instance.downloadStoryVideo(ctx, story)
	}
}

func (instance *Instance) downloadStoryPhoto(ctx Context, story vk.Story) (bool, error) {
	filename := fmt.Sprintf("%s (%d)", story.Date.Format(time.RFC3339), story.ID)
	filename = ctx.DirectoryHandler("stories", filename)
	url := story.Photo.HighestQualityVariantUrl()
	if alreadyDownloaded(filename) {
		instance.Log.Printf("%s; exists", filename)
		return false, nil
	}
	err := instance.download(url, filename)
	if err != nil {
		return false, err
	}
	instance.Log.Printf("%s", filename)
	return true, nil
}

func (instance *Instance) downloadStoryVideo(ctx Context, story vk.Story) (bool, error) {
	filename := fmt.Sprintf("%s (%d)", story.Date.Format(time.RFC3339), story.ID)
	filename = ctx.DirectoryHandler("stories", story.ID)
	if alreadyDownloaded(filename) {
		instance.Log.Printf("%s; exists", filename)
		return false, nil
	}
	url := story.Video.HighestQualityVariantUrl()
	err := instance.download(url, filename)
	if err != nil {
		return false, err
	}
	instance.Log.Printf("%s", filename)
	return true, nil
}

func (instance *Instance) downloadPhotos(ctx Context) error {
	const maxPhotosPerPage = 200
	offset := 0
	for {
		collection, err := instance.Vk.PhotosGetAll(ctx.UserID, maxPhotosPerPage, offset)
		if err != nil {
			ok := handleApiError(err)
			if !ok {
				return err
			}
		}
		if len(collection.Photos) <= 0 {
			break
		}
		for _, photo := range collection.Photos {
			downloaded, err := instance.downloadPhoto(ctx, photo)
			if err != nil {
				return err
			}
			if !downloaded && instance.Config.FastUpdate {
				return nil
			}
			offset += 1
		}
	}
	return nil
}

func (instance *Instance) downloadPhoto(ctx Context, photo vk.Photo) (bool, error) {
	filename := fmt.Sprintf("%s (%d)", photo.Date.Format(time.RFC3339), photo.ID)
	filename = ctx.DirectoryHandler("photos", filename)
	if alreadyDownloaded(filename) {
		instance.Log.Printf("%s; exists", filename)
		return false, nil
	}
	url := photo.HighestQualityVariantUrl()
	err := instance.download(url, filename)
	if err != nil {
		return false, err
	}
	instance.Log.Printf("OK %s", filename)
	return true, nil
}

func (instance *Instance) download(url, filePath string) error {
	fileExtensionCandidate := path.Ext(url)
	fileExtensionParts := strings.Split(fileExtensionCandidate, "?")
	fileExtension := fileExtensionParts[0]
	filePath = filePath + fileExtension
	tmpFilePath := fmt.Sprintf("%s.tmp", filePath)
	output, err := os.Create(tmpFilePath)
	if err != nil {
		return err
	}
	defer output.Close()
	response, err := instance.transport.Get(url)
	if err != nil {
		return err
	}
	defer response.Body.Close()
	_, err = io.Copy(output, response.Body)
	if err != nil {
		return err
	}
	err = os.Rename(tmpFilePath, filePath)
	if err != nil {
		return err
	}
	return nil
}

func handleApiError(err error) bool {
	switch e := err.(type) {
	case vk.InstanceRateLimitedError:
		fmt.Printf("Too many requests; will sleep for %s\n", e.Remaining)
		time.Sleep(e.Remaining)
		return true
	case vk.VkApiError:
		switch e.Code {
		case vk.TooManyRequestsError:
			fmt.Printf("Too many requests; will sleep for %s\n", defaultRateLimitSleepTime)
			time.Sleep(defaultRateLimitSleepTime)
			return true
		case vk.PrivateProfileError:
			fmt.Println("User is private")
			return false
		case vk.UserDeactivatedError:
			fmt.Println("User is deactivated")
			return false
		default:
			fmt.Printf("Unexpected API Error: %s\n", err)
			return false
		}
	default:
		return false
	}
}

func createDirectory(path string) error {
	err := os.MkdirAll(path, os.ModePerm)
	if err != nil {
		return err
	}
	return nil
}

func alreadyDownloaded(path string) bool {
	matches, err := filepath.Glob(path + ".*")
	if err != nil {
		return false
	}
	for _, match := range matches {
		isNotTempFile := !strings.HasSuffix(match, ".tmp")
		return isNotTempFile
	}
	return false
}
