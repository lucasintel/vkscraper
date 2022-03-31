package vk

import "strconv"

type VideoCollection struct {
	Count  int     `json:"count"`
	Videos []Video `json:"items"`
}

type Video struct {
	ID         int              `json:"id"`
	OwnerID    int              `json:"owner_id"`
	AddingDate Timestamp        `json:"adding_date"`
	Variants   VideoVariants    `json:"files"`
	Thumbnails []VideoThumbnail `json:"image"`
}

type VideoVariants struct {
	External  string `json:"external"`
	MP4_2160p string `json:"mp4_2160"`
	MP4_1440p string `json:"mp4_1440"`
	MP4_1080p string `json:"mp4_1080"`
	MP4_720p  string `json:"mp4_720"`
	MP4_480p  string `json:"mp4_480"`
	MP4_360p  string `json:"mp4_360"`
	FLV_320p  string `json:"flv_320"`
	MP4_240p  string `json:"mp4_240"`
	MP4_144p  string `json:"mp4_144"`
}

type VideoThumbnail struct {
	Height int    `json:"height"`
	Width  int    `json:"width"`
	URL    string `json:"url"`
}

func (video Video) HighestQualityThumbnailUrl() string {
	selectedVersion := video.Thumbnails[0]
	for _, version := range video.Thumbnails {
		selectedVersionArea := selectedVersion.Width * selectedVersion.Height
		versionArea := version.Width * version.Height
		if selectedVersionArea < versionArea {
			selectedVersion = version
		}
	}
	return selectedVersion.URL
}

func (video Video) HighestQualityVariantUrl() string {
	findBestVariant := func(values ...string) string {
		for _, value := range values {
			if value == "" {
				continue
			}
			return value
		}
		return ""
	}
	url := findBestVariant(
		video.Variants.MP4_2160p,
		video.Variants.MP4_1440p,
		video.Variants.MP4_1080p,
		video.Variants.MP4_720p,
		video.Variants.MP4_480p,
		video.Variants.MP4_360p,
		video.Variants.FLV_320p,
		video.Variants.MP4_240p,
		video.Variants.MP4_144p,
	)
	return url
}

func (client Client) VideosGet(owner_id, count, offset int) (*VideoCollection, error) {
	params := Params{
		"owner_id": strconv.Itoa(owner_id),
		"count":    strconv.Itoa(count),
		"offset":   strconv.Itoa(offset),
		"extended": "1",
	}
	collection := VideoCollection{}
	err := client.sendMethod("video.get", params, &collection)
	if err != nil {
		return nil, err
	}
	return &collection, nil
}
