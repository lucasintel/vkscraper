package vkscraper

func (instance *Instance) DownloadProfiles(screenNames []string) {
	for index, screenName := range screenNames {
		instance.Log.Printf("[%d] Downloading profile %s\n", index, screenName)
		err := instance.DownloadProfile(screenName)
		if err != nil {
			instance.Log.Printf("[BUG] Could not download profile %s, got error: %s\n", screenName, err)
			continue
		}
	}
}

func (instance *Instance) DownloadProfile(screenName string) error {
	userID, err := instance.Vk.ResolveScreenName(screenName)
	if err != nil {
		ok := handleApiError(err)
		if !ok {
			return err
		}
	}
	return instance.DownloadPhotos(userID)
}
