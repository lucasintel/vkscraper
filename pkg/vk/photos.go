package vk

import (
	"strconv"
)

type PhotoSize struct {
	Height int    `json:"height"`
	Width  int    `json:"width"`
	Type   string `json:"type"`
	URL    string `json:"url"`
}

type PhotoLikes struct {
	UserLikes int `json:"user_likes"`
	Count     int `json:"count"`
}

type PhotoReposts struct {
	Count int `json:"count"`
}

type PhotoComments struct {
	Count int `json:"count"`
}

type PhotoTags struct {
	Count int `json:"count"`
}

type Photo struct {
	ID         int           `json:"id"`
	AlbumID    int           `json:"album_id"`
	OwnerID    int           `json:"owner_id"`
	Text       string        `json:"text"`
	Date       int64         `json:"date"`
	Sizes      []PhotoSize   `json:"sizes"`
	Likes      PhotoLikes    `json:"likes"`
	Reposts    PhotoReposts  `json:"reposts"`
	Comments   PhotoComments `json:"comments"`
	CanComment int           `json:"can_comment"`
	HasTags    bool          `json:"has_tags"`
	Tags       PhotoTags     `json:"tags"`
}

type PhotoCollection struct {
	Count int     `json:"count"`
	Items []Photo `json:"items"`
}

func (client *Client) PhotosGet(owner_id, count, offset int) (*PhotoCollection, error) {
	params := Params{
		"owner_id": strconv.Itoa(owner_id),
		"album_id": "wall",
		"extended": "1",
		"rev":      "1",
		"count":    strconv.Itoa(count),
		"offset":   strconv.Itoa(offset),
	}
	photoCollection := PhotoCollection{}
	err := client.performRequest("photos.get", params, &photoCollection)
	if err != nil {
		return nil, err
	}
	return &photoCollection, nil
}
