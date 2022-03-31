package vk

import "fmt"

type ScreenName struct {
	ID   int    `json:"object_id"`
	Type string `json:"type"`
}

func (client Client) ResolveScreenName(screenName string) (ScreenName, error) {
	params := Params{
		"screen_name": screenName,
	}
	screenNameResponse := ScreenName{}
	err := client.sendMethod("utils.resolveScreenName", params, &screenNameResponse)
	if err != nil {
		return screenNameResponse, fmt.Errorf("could not find user: %s, err: %s", screenName, err)
	}
	return screenNameResponse, nil
}
