package vkscraper

import (
	"fmt"
	"path/filepath"
	"strconv"
	"time"

	"github.com/kandayo/vkscraper/pkg/vk"
)

const (
	maxPhotosPerRequest = 200
	initialOffset       = 0
)

func (instance *Instance) DownloadPhotosByScreenName(screenName string) error {
	userID, err := instance.Vk.ResolveScreenName(screenName)
	if err != nil {
		ok := handleApiError(err)
		if !ok {
			return err
		}
	}
	return instance.DownloadPhotos(userID)
}

func (instance *Instance) DownloadPhotos(userID int) error {
	offset := initialOffset
	for {
		photoCollection, err := instance.Vk.PhotosGet(userID, maxPhotosPerRequest, offset)
		if err != nil {
			ok := handleApiError(err)
			if !ok {
				return err
			}
		}
		if len(photoCollection.Items) <= 0 {
			break
		}
		for _, photo := range photoCollection.Items {
			downloaded, err := instance.downloadPhoto(photo)
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

func (instance *Instance) downloadPhoto(photo vk.Photo) (bool, error) {
	fileName, err := instance.buildPhotoFileName(photo)
	if err != nil {
		return false, err
	}
	instance.Log.Printf("Downloading photo %d.jpg\n", photo.ID)
	alreadyDownloaded := fileExists(fileName + ".jpg")
	if alreadyDownloaded {
		return false, nil
	}
	photoUrl := instance.findHighestQualityPhotoUrl(photo)
	err = instance.download(photoUrl, fileName)
	if err != nil {
		return false, err
	}
	err = saveMetadata(photo, fileName)
	if err != nil {
		return false, err
	}
	return true, nil
}

func (instance *Instance) buildPhotoFileName(photo vk.Photo) (string, error) {
	directory := filepath.Join(instance.Config.BaseDir, strconv.Itoa(photo.OwnerID), "photos")
	err := createDirectory(directory)
	if err != nil {
		return "", err
	}
	postedAt := time.Unix(photo.Date, 0)
	baseFileName := fmt.Sprintf("%s (%d)", postedAt.Format(time.RFC3339), photo.ID)
	fileName := filepath.Join(directory, baseFileName)
	return fileName, nil
}

func (instance *Instance) findHighestQualityPhotoUrl(photo vk.Photo) string {
	selectedVersion := photo.Sizes[0]
	for _, version := range photo.Sizes {
		selectedVersionArea := selectedVersion.Width * selectedVersion.Height
		versionArea := version.Width * version.Height
		if selectedVersionArea < versionArea {
			selectedVersion = version
		}
	}
	return selectedVersion.URL
}
