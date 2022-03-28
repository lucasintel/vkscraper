package vk_test

import (
	"net/http"
	"testing"

	"github.com/jarcoal/httpmock"
	"github.com/kandayo/vkscraper/pkg/vk"
	"github.com/stretchr/testify/assert"
)

func TestPhotosGetError(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	vkResponse := vk.Response{
		Error: vk.ResponseError{
			Code:    vk.TooManyRequestsError,
			Message: "Too many requests per second",
		},
	}
	httpmock.RegisterResponder("GET", "https://api.vk.com/method/photos.get",
		func(req *http.Request) (*http.Response, error) {
			return httpmock.NewJsonResponse(200, vkResponse)
		},
	)

	client := vk.NewClient("token")
	_, err := client.PhotosGet(99, 200, 0)

	expectedError := vk.VkApiError{Code: vk.TooManyRequestsError, Message: "Too many requests per second"}
	assert.Equal(t, err, expectedError)
}

func TestPhotosGetSuccess(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	vkResponse := vk.Response{
		Response: vk.PhotoCollection{
			Count: 10,
			Items: []vk.Photo{},
		},
	}
	httpmock.RegisterResponder("GET", "https://api.vk.com/method/photos.get",
		func(req *http.Request) (*http.Response, error) {
			return httpmock.NewJsonResponse(200, vkResponse)
		},
	)

	client := vk.NewClient("token")
	collection, err := client.PhotosGet(99, 200, 0)

	assert.Nil(t, err)
	assert.Equal(t, collection.Count, 10)
}
