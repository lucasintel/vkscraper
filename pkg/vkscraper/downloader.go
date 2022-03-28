package vkscraper

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path"
	"strings"
	"time"

	"github.com/kandayo/vkscraper/pkg/vk"
)

const (
	userAgent        = "vkscraper/1.0 (+https://github.com/kandayo/vkscraper)"
	timeout          = time.Minute
	defaultSleepTime = 3 * time.Second
)

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

func saveMetadata(resource interface{}, fileName string) error {
	fileName = fileName + ".json"
	file, err := os.OpenFile(fileName, os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	defer file.Close()
	err = json.NewEncoder(file).Encode(&resource)
	if err != nil {
		return err
	}
	return nil
}

func handleApiError(err error) bool {
	switch e := err.(type) {
	case vk.TooManyRequestsClientError:
		message := fmt.Sprintf("Too many requests; will sleep for %s", e.Remaining)
		fmt.Println(message)
		time.Sleep(e.Remaining)
		return true
	case vk.VkApiError:
		if e.Code == vk.TooManyRequestsError {
			message := fmt.Sprintf("Too many requests; will sleep for %s", defaultSleepTime)
			fmt.Println(message)
			time.Sleep(defaultSleepTime)
			return true
		} else {
			message := fmt.Sprintf("Unexpected API Error: %s", err)
			fmt.Println(message)
			return false
		}
	default:
		// message := fmt.Sprintf("[BUG] Unexpected error: %s", err)
		// fmt.Println(message)
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

func fileExists(path string) bool {
	_, err := os.Stat(path)
	if err != nil {
		return false
	}
	return true
}
