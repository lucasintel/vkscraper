package vk_test

import (
	"testing"

	"github.com/kandayo/vkscraper/pkg/vk"
	"github.com/stretchr/testify/assert"
)

func TestNewClient(t *testing.T) {
	t.Parallel()

	client := vk.NewClient("token")

	assert := assert.New(t)
	assert.Equal(client.AccessToken, "token")
	assert.Equal(client.BaseUrl, "api.vk.com")
	assert.Equal(client.Version, "5.131")
	assert.Equal(client.Language, "ru")
	assert.Equal(client.UserAgent, "vkscraper/1.0 (+https://github.com/kandayo/vkscraper)")
}
