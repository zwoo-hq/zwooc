package config

import (
	"encoding/json"
	"os"
)

type Config struct {
	raw      map[string]interface{}
	Profiles []Profile
}

func Load(path string) (Config, error) {
	content, err := os.ReadFile(path)
	if err != nil {
		return Config{}, err
	}

	data := map[string]interface{}{}
	err = json.Unmarshal(content, &data)
	if err != nil {
		return Config{}, err
	}

	return Config{
		raw:      data,
		Profiles: extractProfiles(data),
	}, nil
}

func extractProfiles(data map[string]interface{}) []Profile {
	profiles := []Profile{}

	for projectKey, projectValue := range data {
		if !IsReservedKey(projectKey) {
			project := projectValue.(map[string]interface{})
			for profileKey, profileValue := range project {
				if !IsReservedKey(profileKey) {
					profiles = append(profiles, newProfile(profileValue, profileKey))
				}
			}
		}
	}

	return profiles
}
