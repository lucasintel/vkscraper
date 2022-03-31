package vk

import (
	"strconv"
)

type storiesGetResponse struct {
	Count      int         `json:"count"`
	StoryFeeds []storyFeed `json:"items"`
}

type storyFeed struct {
	Type  string  `json:"type"`
	Items []Story `json:"stories"`
}

type Story struct {
	ID      int       `json:"id"`
	OwnerID int       `json:"owner_id"`
	Date    Timestamp `json:"date"`
	Type    string    `json:"type"`
	Photo   Photo     `json:"photo"`
	Video   Video     `json:"video"`
}

func (story Story) IsPhoto() bool {
	return story.Type == "photo"
}

func (story Story) IsVideo() bool {
	return story.Type == "video"
}

func (client Client) StoriesGet(owner_id int) ([]Story, error) {
	params := Params{
		"owner_id": strconv.Itoa(owner_id),
		"extended": "1",
	}
	collection := storiesGetResponse{}
	err := client.sendMethod("stories.get", params, &collection)
	if err != nil {
		return nil, err
	}
	var stories []Story
	for _, storyFeed := range collection.StoryFeeds {
		stories = append(stories, storyFeed.Items...)
	}
	return stories, nil
}
