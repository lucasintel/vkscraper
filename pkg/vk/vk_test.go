package vk_test

import (
	"testing"

	"github.com/kandayo/vkscraper/pkg/vk"
	"github.com/stretchr/testify/assert"
)

func TestNewClient(t *testing.T) {
	t.Parallel()

	client := vk.NewClient()
	assert.Equal(t, client.AccessToken, "")
	assert.Equal(t, client.BaseUrl, "api.vk.com")
	assert.Equal(t, client.Version, "5.131")
	assert.Equal(t, client.Language, "ru")
	assert.Equal(t, client.UserAgent, "com.vk.vkclient/859 (iPhone; iOS 15.4; Scale/3.00)")
}

func TestClientSetAccessToken(t *testing.T) {
	t.Parallel()

	client := vk.NewClient()
	client.SetAccessToken("token")
	assert.Equal(t, "token", client.AccessToken)
}
