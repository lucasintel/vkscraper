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
	baseUrl   = "api.vk.com"
	version   = "5.131"
	language  = "ru"
	userAgent = "vkscraper/1.0 (+https://github.com/kandayo/vkscraper)"
	timeout   = time.Minute
)

const serviceTokenMaxRequestsPerSecond = 3

type Client struct {
	BaseUrl     string
	AccessToken string
	Version     string
	Language    string
	UserAgent   string
	limiter     *rate.RateLimiter
	transport   *http.Client
}

type ResponseError struct {
	Code    int    `json:"error_code"`
	Message string `json:"error_msg"`
}

type Response struct {
	Error    ResponseError `json:"error"`
	Response interface{}   `json:"response"`
}

type Params map[string]string

func NewClient(accessToken string) *Client {
	limiter := rate.New(serviceTokenMaxRequestsPerSecond, time.Second)
	return &Client{
		BaseUrl:     baseUrl,
		AccessToken: accessToken,
		Version:     version,
		Language:    language,
		UserAgent:   userAgent,
		limiter:     limiter,
		transport: &http.Client{
			Timeout: timeout,
		},
	}
}

func (client *Client) performRequest(method string, params Params, resource interface{}) error {
	requestUrl := client.buildMethodUrl(method, params)
	request, _ := http.NewRequest("GET", requestUrl, nil)
	request.Header.Set("Accept", "application/json; charset=utf-8")
	request.Header.Set("Content-Type", "application/json; charset=utf-8")
	request.Header.Set("User-Agent", client.UserAgent)
	ok, remaining := client.limiter.Try()
	if !ok {
		return TooManyRequestsClientError{Remaining: remaining}
	}
	response, err := client.transport.Do(request)
	if err != nil {
		return err
	}
	defer response.Body.Close()
	if response.StatusCode != http.StatusOK {
		return VkHttpError{Code: response.StatusCode}
	}
	apiResponse := Response{
		Error:    ResponseError{},
		Response: resource,
	}
	err = json.NewDecoder(response.Body).Decode(&apiResponse)
	if err != nil {
		return err
	}
	if apiResponse.Error.Message != "" {
		return VkApiError{
			Code:    apiResponse.Error.Code,
			Message: apiResponse.Error.Message,
		}
	}
	return nil
}

func (client *Client) buildMethodUrl(method string, params Params) string {
	requestPath := fmt.Sprintf("/method/%s", method)
	requestUrl := url.URL{
		Scheme: "https",
		Host:   client.BaseUrl,
		Path:   requestPath,
	}
	query := url.Values{}
	for k, v := range params {
		query.Add(k, v)
	}
	query.Add("access_token", client.AccessToken)
	query.Add("https", "1")
	query.Add("lang", client.Language)
	query.Add("v", client.Version)
	requestUrl.RawQuery = query.Encode()
	return requestUrl.String()
}
