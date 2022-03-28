package vk

type ResolveScreenNameResponse struct {
	ObjectID int    `json:"object_id"`
	Type     string `json:"type"`
}

func (client *Client) ResolveScreenName(screenName string) (int, error) {
	params := Params{
		"screen_name": screenName,
	}
	resolveScreenNameResponse := ResolveScreenNameResponse{}
	err := client.performRequest("utils.resolveScreenName", params, &resolveScreenNameResponse)
	if err != nil {
		return 0, err
	}
	return resolveScreenNameResponse.ObjectID, nil
}
