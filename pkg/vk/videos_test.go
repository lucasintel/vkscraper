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

func TestVideosGetPrivateProfileError(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	fixtureBuffer, _ := ioutil.ReadFile("../../test_fixtures/video.get/private_profile.json")
	fixtureData := string(fixtureBuffer)
	httpmock.RegisterResponder("GET", "https://api.vk.com/method/stories.get",
		httpmock.NewStringResponder(200, fixtureData))

	client := vk.NewClient()
	_, err := client.StoriesGet(99)

	expectedError := vk.VkApiError{
		Code:    vk.PrivateProfileError,
		Message: "This profile is private",
	}
	assert.Equal(t, err, expectedError)
}

func TestVideosGetUserDeactivatedError(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	fixtureBuffer, _ := ioutil.ReadFile("../../test_fixtures/video.get/user_deactivated.json")
	fixtureData := string(fixtureBuffer)
	httpmock.RegisterResponder("GET", "https://api.vk.com/method/stories.get",
		httpmock.NewStringResponder(200, fixtureData))

	client := vk.NewClient()
	_, err := client.StoriesGet(99)

	expectedError := vk.VkApiError{
		Code:    vk.UserDeactivatedError,
		Message: "Access denied: user deactivated",
	}
	assert.Equal(t, err, expectedError)
}

func TestVideosGetFromUser(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	fixtureBuffer, _ := ioutil.ReadFile("../../test_fixtures/video.get/video_from_user.json")
	fixtureData := string(fixtureBuffer)
	httpmock.RegisterResponder("GET", "https://api.vk.com/method/video.get",
		httpmock.NewStringResponder(200, fixtureData))

	client := vk.NewClient()
	collection, err := client.VideosGet(99, 200, 0)

	require.Nil(t, err)
	assert.Equal(t, 270, collection.Count)

	externalVideo := collection.Videos[0]
	assert.Equal(t, 456242027, externalVideo.ID)
	assert.Equal(t, -4908135, externalVideo.OwnerID)
	assert.Equal(t, vk.Timestamp{time.Unix(1643924086, 0)}, externalVideo.AddingDate)
	expectedExternalVideoThumbnails := []vk.VideoThumbnail{
		{URL: "https://static.vk.com/130x96.jpg", Width: 130, Height: 96},
		{URL: "https://static.vk.com/160x120.jpg", Width: 160, Height: 120},
		{URL: "https://static.vk.com/320x240.jpg", Width: 320, Height: 240},
		{URL: "https://static.vk.com/800x450.jpg", Width: 800, Height: 450},
	}
	assert.Equal(t, expectedExternalVideoThumbnails, externalVideo.Thumbnails)
	assert.Equal(t, vk.VideoVariants{External: "https://youtube.com/video.mp4"}, externalVideo.Variants)

	video := collection.Videos[1]
	assert.Equal(t, 456239469, video.ID)
	assert.Equal(t, 15594498, video.OwnerID)
	assert.Equal(t, vk.Timestamp{time.Unix(1640165512, 0)}, video.AddingDate)
	expectedVideoThumbnails := []vk.VideoThumbnail{
		{URL: "https://static.vk.com/file.jpg", Width: 130, Height: 96},
		{URL: "https://static.vk.com/file.jpg", Width: 160, Height: 120},
		{URL: "https://static.vk.com/file.jpg", Width: 320, Height: 240},
		{URL: "https://static.vk.com/file.jpg", Width: 800, Height: 450},
		{URL: "https://static.vk.com/file.jpg", Width: 720, Height: 720},
		{URL: "https://static.vk.com/file.jpg", Width: 320, Height: 320},
		{URL: "https://static.vk.com/file.jpg", Width: 720, Height: 720},
		{URL: "https://static.vk.com/file.jpg", Width: 1024, Height: 1024},
		{URL: "https://static.vk.com/file.jpg", Width: 4096, Height: 4096},
	}
	assert.Equal(t, expectedVideoThumbnails, video.Thumbnails)
	expectedVideoVariants := vk.VideoVariants{
		MP4_144p:  "https://static.vk.com/144.mp4",
		MP4_240p:  "https://static.vk.com/240.mp4",
		MP4_360p:  "https://static.vk.com/360.mp4",
		MP4_480p:  "https://static.vk.com/480.mp4",
		MP4_720p:  "https://static.vk.com/720.mp4",
		MP4_1080p: "https://static.vk.com/1080.mp4",
	}
	assert.Equal(t, expectedVideoVariants, video.Variants)
}

func TestVideosGetFromCommunity(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	fixtureBuffer, _ := ioutil.ReadFile("../../test_fixtures/video.get/video_from_community.json")
	fixtureData := string(fixtureBuffer)
	httpmock.RegisterResponder("GET", "https://api.vk.com/method/video.get",
		httpmock.NewStringResponder(200, fixtureData))

	client := vk.NewClient()
	collection, err := client.VideosGet(99, 200, 0)

	require.Nil(t, err)
	assert.Equal(t, 51, collection.Count)
	video := collection.Videos[0]
	assert.Equal(t, 456239308, video.ID)
	assert.Equal(t, -35005, video.OwnerID)
	assert.Equal(t, vk.Timestamp{time.Unix(1644932218, 0)}, video.AddingDate)
	expectedExternalVideoThumbnails := []vk.VideoThumbnail{
		{URL: "https://static.vk.com/file.jpg", Width: 130, Height: 96},
		{URL: "https://static.vk.com/file.jpg", Width: 160, Height: 120},
		{URL: "https://static.vk.com/file.jpg", Width: 320, Height: 240},
		{URL: "https://static.vk.com/file.jpg", Width: 800, Height: 450},
		{URL: "https://static.vk.com/file.jpg", Width: 1280, Height: 720},
		{URL: "https://static.vk.com/file.jpg", Width: 320, Height: 180},
		{URL: "https://static.vk.com/file.jpg", Width: 720, Height: 405},
		{URL: "https://static.vk.com/file.jpg", Width: 1024, Height: 576},
		{URL: "https://static.vk.com/file.jpg", Width: 4096, Height: 2304},
	}
	assert.Equal(t, expectedExternalVideoThumbnails, video.Thumbnails)
	expectedVideoVariants := vk.VideoVariants{
		MP4_144p:  "https://static.vk.com/144.mp4",
		MP4_240p:  "https://static.vk.com/240.mp4",
		MP4_360p:  "https://static.vk.com/360.mp4",
		MP4_480p:  "https://static.vk.com/480.mp4",
		MP4_720p:  "https://static.vk.com/720.mp4",
		MP4_1080p: "https://static.vk.com/1080.mp4",
	}
	assert.Equal(t, expectedVideoVariants, video.Variants)
}
