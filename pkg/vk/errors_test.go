package vk_test

import (
	"net/http"
	"testing"
	"time"

	"github.com/kandayo/vkscraper/pkg/vk"
	"github.com/stretchr/testify/assert"
)

func TestVkHttpErrorMessage(t *testing.T) {
	t.Parallel()
	e := vk.VkHttpError{Code: http.StatusNotFound}
	assert.Equal(t, e.Error(), "VK responded with status 404")
}

func TestVkApiErrorMessage(t *testing.T) {
	t.Parallel()
	e := vk.VkApiError{
		Code:    vk.TooManyRequestsError,
		Message: "Too many requests per second",
	}
	assert.Equal(t, e.Error(), "Too many requests per second")
}

func TestInstanceRateLimitedErrorMessage(t *testing.T) {
	t.Parallel()
	remaining := 10 * time.Second
	e := vk.InstanceRateLimitedError{
		Remaining: remaining,
	}
	assert.Equal(t, e.Error(), "Too many requests; please wait 10s")
}
