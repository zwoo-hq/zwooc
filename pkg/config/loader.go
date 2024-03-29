package config

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"github.com/zwoo-hq/zwooc/pkg/model"
)

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

	c := Config{
		baseDir: filepath.Dir(path),
		raw:     data,
	}
	err = c.init()
	return c, err
}

func (c *Config) init() error {
	var err error
	c.profiles, err = c.loadProfiles()
	if err != nil {
		return err
	}

	c.fragments, err = c.loadFragments()
	if err != nil {
		return err
	}

	c.compounds, err = c.loadCompounds()
	return err
}

func (c Config) GetProfiles() []Profile {
	return c.profiles
}

func (c Config) loadProfiles() ([]Profile, error) {
	profiles := []Profile{}

	for projectKey, projectValue := range c.raw {
		if !IsReservedKey(projectKey) {
			project := projectValue.(map[string]interface{})
			var projectAdapter string
			if adapter, ok := project[model.KeyAdapter]; ok {
				projectAdapter = adapter.(string)
			} else {
				return []Profile{}, fmt.Errorf("project '%s' is missing adapter", projectKey)
			}
			projectDirectory := projectKey
			if directory, ok := project[model.KeyDirectory]; ok {
				projectDirectory = directory.(string)
			}

			for profileKey, profileValue := range project {
				if !IsReservedKey(profileKey) {
					newProfile := Profile{
						name:      profileKey,
						adapter:   projectAdapter,
						directory: filepath.Join(c.baseDir, projectDirectory),
						raw:       profileValue.(map[string]interface{}),
					}
					profiles = append(profiles, newProfile)
				}
			}
		}
	}

	return profiles, nil
}

func (c Config) GetFragments() []Fragment {
	return c.fragments
}

func (c Config) loadFragments() ([]Fragment, error) {
	fragments := []Fragment{}

	for projectKey, projectValue := range c.raw {
		if !IsReservedKey(projectKey) {
			project := projectValue.(map[string]interface{})
			if fragmentDefinitions, ok := project[model.KeyFragment]; ok {
				for fragmentKey, fragmentValue := range fragmentDefinitions.(map[string]interface{}) {
					projectDirectory := projectKey
					if directory, ok := project[model.KeyDirectory]; ok {
						projectDirectory = directory.(string)
					}

					newFragment := Fragment{
						name:      fragmentKey,
						directory: filepath.Join(c.baseDir, projectDirectory),
						raw:       fragmentValue,
					}
					fragments = append(fragments, newFragment)
				}
			}
		}
	}

	if fragmentDefinitions, ok := c.raw[model.KeyFragment]; ok {
		for fragmentKey, fragmentValue := range fragmentDefinitions.(map[string]interface{}) {
			newFragment := Fragment{
				name:      fragmentKey,
				directory: c.baseDir,
				raw:       fragmentValue,
			}
			fragments = append(fragments, newFragment)
		}
	}

	return fragments, nil
}

func (c Config) GetCompounds() []Compound {
	return c.compounds
}

func (c Config) loadCompounds() ([]Compound, error) {
	compounds := []Compound{}

	if compoundDefinitions, ok := c.raw[model.KeyCompound]; ok {
		for compoundKey, compoundValue := range compoundDefinitions.(map[string]interface{}) {
			newCompound := Compound{
				name:      compoundKey,
				directory: c.baseDir,
				raw:       compoundValue.(map[string]interface{}),
			}
			compounds = append(compounds, newCompound)
		}
	}
	return compounds, nil
}
