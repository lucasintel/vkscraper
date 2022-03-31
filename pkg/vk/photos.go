package vk

import (
	"strconv"
)

type PhotoCollection struct {
	Count  int     `json:"count"`
	Photos []Photo `json:"items"`
}

type Photo struct {
	ID       int            `json:"id"`
	OwnerID  int            `json:"owner_id"`
	Date     Timestamp      `json:"date"`
	Variants []PhotoVariant `json:"sizes"`
}

type PhotoVariant struct {
	Height int    `json:"height"`
	Width  int    `json:"width"`
	URL    string `json:"url"`
}

func (photo Photo) HighestQualityVariantUrl() string {
	selectedVersion := photo.Variants[0]
	for _, version := range photo.Variants {
		selectedVersionArea := selectedVersion.Width * selectedVersion.Height
		versionArea := version.Width * version.Height
		if selectedVersionArea < versionArea {
			selectedVersion = version
		}
	}
	return selectedVersion.URL
}

func (client Client) PhotosGetAll(owner_id, count, offset int) (PhotoCollection, error) {
	params := Params{
		"owner_id":          strconv.Itoa(owner_id),
		"extended":          "1",
		"offset":            strconv.Itoa(offset),
		"count":             strconv.Itoa(count),
		"photo_sizes":       "1",
		"no_service_albums": "0",
		"need_hidden":       "1",
		"skip_hidden":       "0",
	}
	photoCollection := PhotoCollection{}
	err := client.sendMethod("photos.getAll", params, &photoCollection)
	if err != nil {
		return photoCollection, err
	}
	return photoCollection, nil
}
