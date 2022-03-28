package vk_test

import (
	"net/http"
	"testing"

	"github.com/jarcoal/httpmock"
	"github.com/kandayo/vkscraper/pkg/vk"
	"github.com/stretchr/testify/assert"
)

func TestResolveScreenNameError(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	vkResponse := vk.Response{
		Error: vk.ResponseError{
			Code:    vk.TooManyRequestsError,
			Message: "Too many requests per second",
		},
	}
	httpmock.RegisterResponder("GET", "https://api.vk.com/method/utils.resolveScreenName",
		func(req *http.Request) (*http.Response, error) {
			return httpmock.NewJsonResponse(200, vkResponse)
		},
	)

	client := vk.NewClient("token")
	_, err := client.ResolveScreenName("charlidamelio")

	expectedError := vk.VkApiError{Code: vk.TooManyRequestsError, Message: "Too many requests per second"}
	assert.Equal(t, err, expectedError)
}

func TestResolveScreenNameSuccess(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	vkResponse := vk.Response{
		Response: vk.ResolveScreenNameResponse{
			ObjectID: 99,
			Type:     "user",
		},
	}
	httpmock.RegisterResponder("GET", "https://api.vk.com/method/utils.resolveScreenName",
		func(req *http.Request) (*http.Response, error) {
			return httpmock.NewJsonResponse(200, vkResponse)
		},
	)

	client := vk.NewClient("token")
	userID, err := client.ResolveScreenName("charlidamelio")

	assert.Nil(t, err)
	assert.Equal(t, userID, 99)
}
