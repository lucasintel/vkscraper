package vk_test

import (
	"io/ioutil"
	"testing"
	"time"

	"github.com/jarcoal/httpmock"
	"github.com/kandayo/vkscraper/pkg/vk"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestStoriesGetError(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	fixtureBuffer, _ := ioutil.ReadFile("../../test_fixtures/stories.get/invalid_params_error.json")
	fixtureData := string(fixtureBuffer)
	httpmock.RegisterResponder("GET", "https://api.vk.com/method/stories.get",
		httpmock.NewStringResponder(200, fixtureData))

	client := vk.NewClient()
	_, err := client.StoriesGet(99)

	expectedError := vk.VkApiError{
		Code:    vk.InvalidParamsError,
		Message: "One of the parameters specified was missing or invalid: owner_id not integer",
	}
	assert.Equal(t, err, expectedError)
}

func TestStoriesGetNoStories(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	fixtureBuffer, _ := ioutil.ReadFile("../../test_fixtures/stories.get/no_stories.json")
	fixtureData := string(fixtureBuffer)
	httpmock.RegisterResponder("GET", "https://api.vk.com/method/stories.get",
		httpmock.NewStringResponder(200, fixtureData))

	client := vk.NewClient()
	stories, err := client.StoriesGet(99)

	require.Nil(t, err)
	assert.Len(t, stories, 0)
}

func TestStoriesGetFromUser(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	fixtureBuffer, _ := ioutil.ReadFile("../../test_fixtures/stories.get/story_from_user.json")
	fixtureData := string(fixtureBuffer)
	httpmock.RegisterResponder("GET", "https://api.vk.com/method/stories.get",
		httpmock.NewStringResponder(200, fixtureData))

	client := vk.NewClient()
	stories, err := client.StoriesGet(10)

	require.Nil(t, err)
	assert.Len(t, stories, 8)

	videoStory := stories[0]
	assert.Equal(t, 456239207, videoStory.ID)
	assert.Equal(t, 314833033, videoStory.OwnerID)
	assert.Equal(t, vk.Timestamp{time.Unix(1648552720, 0)}, videoStory.Date)
	assert.Equal(t, "video", videoStory.Type)
	assert.True(t, videoStory.IsVideo())
	assert.False(t, videoStory.IsPhoto())
	assert.Equal(t, "https://static.vk.com/800x450.jpg", videoStory.Video.HighestQualityThumbnailUrl())
	assert.Equal(t, "https://static.vk.com/720p.mp4", videoStory.Video.HighestQualityVariantUrl())

	photoStory := stories[7]
	assert.Equal(t, 456239214, photoStory.ID)
	assert.Equal(t, 314833033, photoStory.OwnerID)
	assert.Equal(t, vk.Timestamp{time.Unix(1648568170, 0)}, photoStory.Date)
	assert.Equal(t, "photo", photoStory.Type)
	assert.False(t, photoStory.IsVideo())
	assert.True(t, photoStory.IsPhoto())
	assert.Equal(t, "https://static.vk.com/1920x1080.jpg", photoStory.Photo.HighestQualityVariantUrl())
}

func TestStoriesGetFromCommunity(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	fixtureBuffer, _ := ioutil.ReadFile("../../test_fixtures/stories.get/story_from_community.json")
	fixtureData := string(fixtureBuffer)
	httpmock.RegisterResponder("GET", "https://api.vk.com/method/stories.get",
		httpmock.NewStringResponder(200, fixtureData))

	client := vk.NewClient()
	stories, err := client.StoriesGet(10)

	require.Nil(t, err)
	assert.Len(t, stories, 2)

	story := stories[0]
	assert.Equal(t, 456239124, story.ID)
	assert.Equal(t, -43970955, story.OwnerID)
	assert.Equal(t, vk.Timestamp{time.Unix(1648549795, 0)}, story.Date)
	assert.Equal(t, "photo", story.Type)
	assert.False(t, story.IsVideo())
	assert.True(t, story.IsPhoto())
	assert.Equal(t, "https://static.vk.com/1920x1080.jpg", story.Photo.HighestQualityVariantUrl())

	story = stories[1]
	assert.Equal(t, 456239125, story.ID)
	assert.Equal(t, -43970955, story.OwnerID)
	assert.Equal(t, vk.Timestamp{time.Unix(1648575106, 0)}, story.Date)
	assert.Equal(t, "photo", story.Type)
	assert.False(t, story.IsVideo())
	assert.True(t, story.IsPhoto())
	assert.Equal(t, "https://static.vk.com/1920x1079.jpg", story.Photo.HighestQualityVariantUrl())
}
