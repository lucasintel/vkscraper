package vk_test

import (
	"io/ioutil"
	"testing"

	"github.com/jarcoal/httpmock"
	"github.com/kandayo/vkscraper/pkg/vk"
	"github.com/stretchr/testify/assert"
)

func TestResolveScreenNameNotFound(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	fixtureBuffer, _ := ioutil.ReadFile("../../test_fixtures/utils.resolveScreenName/not_found.json")
	fixtureData := string(fixtureBuffer)
	httpmock.RegisterResponder("GET", "https://api.vk.com/method/utils.resolveScreenName",
		httpmock.NewStringResponder(200, fixtureData))

	client := vk.NewClient()
	_, err := client.ResolveScreenName("klavdiacoca")

	assert.NotNil(t, err)
}

func TestResolveScreenNameSuccess(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	fixtureBuffer, _ := ioutil.ReadFile("../../test_fixtures/utils.resolveScreenName/success.json")
	fixtureData := string(fixtureBuffer)
	httpmock.RegisterResponder("GET", "https://api.vk.com/method/utils.resolveScreenName",
		httpmock.NewStringResponder(200, fixtureData))

	client := vk.NewClient()
	screenName, err := client.ResolveScreenName("klavdiacoca")

	assert.Nil(t, err)
	assert.Equal(t, 15594498, screenName.ID)
	assert.Equal(t, "user", screenName.Type)
}
