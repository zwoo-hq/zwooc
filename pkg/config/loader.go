package config

import (
	"encoding/json"
	"fmt"
	"os"
)

type Config struct {
	raw map[string]interface{}
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
		raw: data,
	}, nil
}

func (c Config) GetProfiles() ([]Profile, error) {
	profiles := []Profile{}

	for projectKey, projectValue := range c.raw {
		if !IsReservedKey(projectKey) {
			project := projectValue.(map[string]interface{})
			var projectAdapter string
			if adapter, ok := project[KeyAdapter]; ok {
				projectAdapter = adapter.(string)
			} else {
				return []Profile{}, fmt.Errorf("project '%s' is missing adapter", projectKey)
			}

			for profileKey, profileValue := range project {
				if !IsReservedKey(profileKey) {
					newProfile := newProfile(profileKey, projectAdapter, projectKey, profileValue.(map[string]interface{}))
					profiles = append(profiles, newProfile)
				}
			}
		}
	}

	return profiles, nil
}

func (c Config) GetFragments() ([]Fragment, error) {
	fragments := []Fragment{}

	for projectKey, projectValue := range c.raw {
		if !IsReservedKey(projectKey) {
			project := projectValue.(map[string]interface{})
			if fragmentDefinitions, ok := project[KeyFragment]; ok {
				for fragmentKey, fragmentValue := range fragmentDefinitions.(map[string]interface{}) {
					newFragment := newFragment(fragmentKey, fragmentValue.(map[string]interface{}))
					fragments = append(fragments, newFragment)
				}
			}
		}
	}

	if fragmentDefinitions, ok := c.raw[KeyFragment]; ok {
		for fragmentKey, fragmentValue := range fragmentDefinitions.(map[string]interface{}) {
			newFragment := newFragment(fragmentKey, fragmentValue.(map[string]interface{}))
			fragments = append(fragments, newFragment)
		}
	}

	return fragments, nil
}

func (c Config) GetCompounds() ([]Compound, error) {
	compounds := []Compound{}

	if compoundDefinitions, ok := c.raw[KeyCompound]; ok {
		for compoundKey, compoundValue := range compoundDefinitions.(map[string]interface{}) {
			newCompound := newCompound(compoundKey, compoundValue.(map[string]interface{}))
			compounds = append(compounds, newCompound)
		}
	}
	return compounds, nil
}
