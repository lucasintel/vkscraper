package vk

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"time"

	"github.com/beefsack/go-rate"
)

const (
	baseUrl              = "api.vk.com"
	version              = "5.131"
	language             = "ru"
	userAgent            = "com.vk.vkclient/859 (iPhone; iOS 15.4; Scale/3.00)"
	timeout              = time.Minute
	maxRequestsPerSecond = 3
)

type Client struct {
	AccessToken string
	BaseUrl     string
	Version     string
	Language    string
	UserAgent   string
	Limiter     RateLimiter
	Transport   *http.Client
}

type RateLimiter interface {
	Try() (bool, time.Duration)
}

type MethodError struct {
	Code    int    `json:"error_code"`
	Message string `json:"error_msg"`
}

type MethodResponse struct {
	Error    MethodError `json:"error"`
	Response interface{} `json:"response"`
}

type Params map[string]string

func NewClient() *Client {
	limiter := rate.New(maxRequestsPerSecond, time.Second)
	transport := &http.Client{
		Timeout: timeout,
	}
	return &Client{
		BaseUrl:   baseUrl,
		Version:   version,
		Language:  language,
		UserAgent: userAgent,
		Limiter:   limiter,
		Transport: transport,
	}
}

func (client *Client) SetAccessToken(accessToken string) {
	client.AccessToken = accessToken
}

func (client Client) sendMethod(method string, params Params, resource interface{}) error {
	url := client.buildMethodUrl(method, params)
	response, err := client.doHttp(url)
	if err != nil {
		return err
	}
	defer response.Body.Close()
	if response.StatusCode != http.StatusOK {
		return VkHttpError{Code: response.StatusCode}
	}
	result := MethodResponse{
		Error:    MethodError{},
		Response: resource,
	}
	err = json.NewDecoder(response.Body).Decode(&result)
	if err != nil {
		return err
	}
	if (result.Error != MethodError{}) {
		return VkApiError{
			Code:    result.Error.Code,
			Message: result.Error.Message,
		}
	}
	return nil
}

func (client Client) buildMethodUrl(method string, queryParams Params) string {
	requestPath := fmt.Sprintf("/method/%s", method)
	requestUrl := url.URL{Scheme: "https", Host: client.BaseUrl, Path: requestPath}
	query := url.Values{}
	for param, value := range queryParams {
		query.Add(param, value)
	}
	query.Add("access_token", client.AccessToken)
	query.Add("https", "1")
	query.Add("lang", client.Language)
	query.Add("v", client.Version)
	requestUrl.RawQuery = query.Encode()
	return requestUrl.String()
}

func (client Client) doHttp(url string) (*http.Response, error) {
	request, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	request.Header.Set("Accept", "application/json; charset=utf-8")
	request.Header.Set("Content-Type", "application/json; charset=utf-8")
	request.Header.Set("User-Agent", client.UserAgent)
	ok, remaining := client.Limiter.Try()
	if !ok {
		return nil, InstanceRateLimitedError{Remaining: remaining}
	}
	return client.Transport.Do(request)
}
