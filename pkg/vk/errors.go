package vk

import (
	"fmt"
	"time"
)

// VK error codes as specified on documentation.
// See: https://vk.com/dev/errors
const (
	UnknownError                                          = 1
	ApplicationDisabledError                              = 2
	UnknownMethodError                                    = 3
	IncorrectSignatureError                               = 4
	UserAuthorizationFailedError                          = 5
	TooManyRequestsError                                  = 6
	PermissionDeniedError                                 = 7
	InvalidRequestError                                   = 8
	FloodControlError                                     = 9
	InternalServerError                                   = 10
	TestModeError                                         = 11
	CaptchaNeededError                                    = 14
	AccessDeniedError                                     = 15
	HTTPAuthorizationFailedError                          = 16
	ValidationRequiredError                               = 17
	UserDeletedOrBannedError                              = 18
	PermissionDeniedForNonStandaloneAppsError             = 20
	PermissionAllowedOnlyForStandaloneAndOpenAPIAppsError = 21
	MethodDisabledError                                   = 23
	ConfirmationRequiredError                             = 24
	GroupAuthorizationFailedError                         = 27
	ApplicationAuthorizationFailedError                   = 28
	RateLimitReachedError                                 = 29
	PrivateProfileError                                   = 30
	NotImplementedError                                   = 33
	ClientVersionDeprecatedError                          = 34
	UserBannedError                                       = 37
	UnknownApplicationError                               = 38
	UnknownUserError                                      = 39
	UnknownGroupError                                     = 40
	AdditionalSignupRequired                              = 41
	IPNotAllowedError                                     = 42
	MissingOrInvalidParametersError                       = 100
	InvalidApiIdError                                     = 101
	InvalidUserIdError                                    = 113
	InvalidTimestampError                                 = 150
	AlbumAccessDeniedError                                = 200
	AudioAccessDeniedError                                = 201
	GroupAccessDeniedError                                = 203
	FullAlbumError                                        = 300
	VotesProcessingDisabledError                          = 500
	NoAccessToOperationSpecifiedError                     = 600
	AdsError                                              = 603
	AnonymousTokenExpiredError                            = 1114
	AnonymousTokenInvalidError                            = 1116
	RecaptchaNeededError                                  = 3300
	PhoneValidationNeededError                            = 3301
	PasswordValidationNeededError                         = 3302
	OtpAppValidationNeededError                           = 3303
	EmailConfirmationNeededError                          = 3304
	AssertVotesError                                      = 3605
	TokenExtensionRequiredError                           = 3609
	UserDeactivatedError                                  = 3610
	ServiceDeactivatedForUserError                        = 3611
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

type TooManyRequestsClientError struct {
	Remaining time.Duration
}

func (e TooManyRequestsClientError) Error() string {
	return fmt.Sprintf("Too many requests; please wait %s", e.Remaining)
}
