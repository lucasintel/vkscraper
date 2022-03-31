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

func TestPhotosGetAllPrivateProfileError(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	fixtureBuffer, _ := ioutil.ReadFile("../../test_fixtures/photos.getAll/private_profile.json")
	fixtureData := string(fixtureBuffer)
	httpmock.RegisterResponder("GET", "https://api.vk.com/method/photos.getAll",
		httpmock.NewStringResponder(200, fixtureData))

	client := vk.NewClient()
	_, err := client.PhotosGetAll(99, 200, 0)

	expectedError := vk.VkApiError{
		Code:    vk.PrivateProfileError,
		Message: "This profile is private",
	}
	assert.Equal(t, err, expectedError)
}

func TestPhotosGetAllFromUser(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	fixtureBuffer, _ := ioutil.ReadFile("../../test_fixtures/photos.getAll/photo_from_user.json")
	fixtureData := string(fixtureBuffer)
	httpmock.RegisterResponder("GET", "https://api.vk.com/method/photos.getAll",
		httpmock.NewStringResponder(200, fixtureData))

	client := vk.NewClient()
	collection, err := client.PhotosGetAll(99, 200, 0)

	require.Nil(t, err)
	assert.Equal(t, 2447, collection.Count)
	photo := collection.Photos[0]
	assert.Equal(t, 457244728, photo.ID)
	assert.Equal(t, 15594498, photo.OwnerID)
	assert.Equal(t, vk.Timestamp{time.Unix(1648556794, 0)}, photo.Date)
	expectedVariants := []vk.PhotoVariant{
		{URL: "https://static.vk.com/file.jpg", Height: 130, Width: 104},
		{URL: "https://static.vk.com/file.jpg", Height: 163, Width: 130},
		{URL: "https://static.vk.com/file.jpg", Height: 250, Width: 200},
		{URL: "https://static.vk.com/file.jpg", Height: 400, Width: 320},
		{URL: "https://static.vk.com/file.jpg", Height: 638, Width: 510},
		{URL: "https://static.vk.com/file.jpg", Height: 75, Width: 60},
		{URL: "https://static.vk.com/file.jpg", Height: 1688, Width: 1350},
		{URL: "https://static.vk.com/file.jpg", Height: 604, Width: 483},
		{URL: "https://static.vk.com/file.jpg", Height: 807, Width: 646},
		{URL: "https://static.vk.com/file.jpg", Height: 1080, Width: 864},
	}
	assert.Equal(t, expectedVariants, photo.Variants)
}

func TestPhotosGetAllFromCommunity(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	fixtureBuffer, _ := ioutil.ReadFile("../../test_fixtures/photos.getAll/photo_from_community.json")
	fixtureData := string(fixtureBuffer)
	httpmock.RegisterResponder("GET", "https://api.vk.com/method/photos.getAll",
		httpmock.NewStringResponder(200, fixtureData))

	client := vk.NewClient()
	collection, err := client.PhotosGetAll(99, 200, 0)

	require.Nil(t, err)
	assert.Equal(t, 1793, collection.Count)
	photo := collection.Photos[0]
	assert.Equal(t, 457242502, photo.ID)
	assert.Equal(t, -43970955, photo.OwnerID)
	assert.Equal(t, vk.Timestamp{time.Unix(1648412246, 0)}, photo.Date)
	expectedVariants := []vk.PhotoVariant{
		{URL: "https://static.vk.com/file.jpg", Height: 75, Width: 56},
		{URL: "https://static.vk.com/file.jpg", Height: 130, Width: 97},
		{URL: "https://static.vk.com/file.jpg", Height: 604, Width: 453},
		{URL: "https://static.vk.com/file.jpg", Height: 807, Width: 605},
		{URL: "https://static.vk.com/file.jpg", Height: 1080, Width: 810},
		{URL: "https://static.vk.com/file.jpg", Height: 1600, Width: 1200},
		{URL: "https://static.vk.com/file.jpg", Height: 173, Width: 130},
		{URL: "https://static.vk.com/file.jpg", Height: 267, Width: 200},
		{URL: "https://static.vk.com/file.jpg", Height: 427, Width: 320},
		{URL: "https://static.vk.com/file.jpg", Height: 680, Width: 510},
	}
	assert.Equal(t, expectedVariants, photo.Variants)
}
