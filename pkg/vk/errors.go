package vk

import (
	"fmt"
	"time"
)

// VK error codes as specified on documentation.
// See: https://vk.com/dev/errors
const (
	TooManyRequestsError = 6
	UserDeactivatedError = 15
	PrivateProfileError  = 30
	InvalidParamsError   = 100
)

type VkHttpError struct {
	Code int
}

func (e VkHttpError) Error() string {
	return fmt.Sprintf("VK responded with status %d", e.Code)
}

type VkApiError struct {
	Code    int
	Message string
}

func (e VkApiError) Error() string {
	return e.Message
}

type InstanceRateLimitedError struct {
	Remaining time.Duration
}

func (e InstanceRateLimitedError) Error() string {
	return fmt.Sprintf("Too many requests; please wait %s", e.Remaining)
}
