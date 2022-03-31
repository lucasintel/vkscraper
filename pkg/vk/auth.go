package vk

import (
	"encoding/json"
	"net/http"
	"net/url"
)

const (
	oauthUrl           = "https://oauth.vk.com/token"
	iPhoneClientID     = "3140623"
	iPhoneClientSecret = "VeWdmVclDCtn6ihuP1nt"
)

type UserToken struct {
	AccessToken string `json:"access_token"`
	ExpiresIn   int    `json:"expires_in"`
	UserID      int    `json:"user_id"`
}

type AuthError struct {
	Code        string `json:"error"`
	Description string `json:"error_description"`
}

func (e AuthError) Error() string {
	return e.Description
}

func (client *Client) Login(username, password string) error {
	params := url.Values{}
	params.Add("client_id", iPhoneClientID)
	params.Add("client_secret", iPhoneClientSecret)
	params.Add("grant_type", "password")
	params.Add("scope", "friends,photos,audio,video,stories,pages,status,notes,wall,offline,docs,groups,stats,email")
	params.Add("username", username)
	params.Add("password", password)
	url, _ := url.Parse(oauthUrl)
	url.RawQuery = params.Encode()
	requestUrl := url.String()
	response, err := client.doHttp(requestUrl)
	if err != nil {
		return err
	}
	defer response.Body.Close()
	if response.StatusCode != http.StatusOK {
		authError := AuthError{}
		err = json.NewDecoder(response.Body).Decode(&authError)
		if err != nil {
			return err
		}
		return authError
	}
	userToken := UserToken{}
	err = json.NewDecoder(response.Body).Decode(&userToken)
	if err != nil {
		return err
	}
	client.SetAccessToken(userToken.AccessToken)
	return nil
}
