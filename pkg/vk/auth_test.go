package vk_test

import (
	"io/ioutil"
	"testing"

	"github.com/jarcoal/httpmock"
	"github.com/kandayo/vkscraper/pkg/vk"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestAuthErrorMessage(t *testing.T) {
	t.Parallel()
	e := vk.AuthError{Description: "client_secret is incorrect"}
	assert.Equal(t, "client_secret is incorrect", e.Error())
}

func TestAuthInvalidRequestError(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	fixtureBuffer, _ := ioutil.ReadFile("../../test_fixtures/auth/invalid_request_error.json")
	fixtureData := string(fixtureBuffer)
	httpmock.RegisterResponder("GET", "https://oauth.vk.com/token?client_id=3140623&client_secret=VeWdmVclDCtn6ihuP1nt&grant_type=password&password=pass&scope=friends%2Cphotos%2Caudio%2Cvideo%2Cstories%2Cpages%2Cstatus%2Cnotes%2Cwall%2Coffline%2Cdocs%2Cgroups%2Cstats%2Cemail&username=user",
		httpmock.NewStringResponder(401, fixtureData))

	client := vk.NewClient()
	err := client.Login("user", "pass")

	expectedError := vk.AuthError{
		Code:        "invalid_request",
		Description: "Too many invalid requests per 15 seconds",
	}
	require.NotNil(t, err)
	assert.Equal(t, expectedError, err)
}

func TestAuthInvalidClientError(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	fixtureBuffer, _ := ioutil.ReadFile("../../test_fixtures/auth/invalid_client_error.json")
	fixtureData := string(fixtureBuffer)
	httpmock.RegisterResponder("GET", "https://oauth.vk.com/token?client_id=3140623&client_secret=VeWdmVclDCtn6ihuP1nt&grant_type=password&password=pass&scope=friends%2Cphotos%2Caudio%2Cvideo%2Cstories%2Cpages%2Cstatus%2Cnotes%2Cwall%2Coffline%2Cdocs%2Cgroups%2Cstats%2Cemail&username=user",
		httpmock.NewStringResponder(401, fixtureData))

	client := vk.NewClient()
	err := client.Login("user", "pass")

	expectedError := vk.AuthError{
		Code:        "invalid_client",
		Description: "client_secret is incorrect",
	}
	require.NotNil(t, err)
	assert.Equal(t, expectedError, err)
}

func TestAuthSuccess(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	fixtureBuffer, _ := ioutil.ReadFile("../../test_fixtures/auth/success.json")
	fixtureData := string(fixtureBuffer)
	httpmock.RegisterResponder("GET", "https://oauth.vk.com/token?client_id=3140623&client_secret=VeWdmVclDCtn6ihuP1nt&grant_type=password&password=pass&scope=friends%2Cphotos%2Caudio%2Cvideo%2Cstories%2Cpages%2Cstatus%2Cnotes%2Cwall%2Coffline%2Cdocs%2Cgroups%2Cstats%2Cemail&username=user",
		httpmock.NewStringResponder(200, fixtureData))

	client := vk.NewClient()
	err := client.Login("user", "pass")

	require.Nil(t, err)
	assert.Equal(t, "access_token", client.AccessToken)
}
