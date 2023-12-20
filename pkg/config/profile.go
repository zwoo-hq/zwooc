package config

import "fmt"

type Profile struct {
	name      string
	adapter   string
	directory string
	raw       map[string]interface{}
}

func newProfile(name, adapter, directory string, data map[string]interface{}) Profile {
	p := Profile{
		name:      name,
		adapter:   adapter,
		directory: directory,
		raw:       data,
	}
	return p
}

func (p Profile) Name() string {
	return p.name
}

func (p Profile) GetConfig(mode string) (RunConfig, error) {
	if !IsValidRunMode(mode) {
		return RunConfig{}, fmt.Errorf("invalid run mode: %s", mode)
	}

	if options, ok := p.raw[mode]; ok {
		config := RunConfig{
			Name:      p.name,
			Adapter:   p.adapter,
			Directory: p.directory,
		}

		if optionsMap, ok := options.(map[string]interface{}); ok {
			config.Options = optionsMap
		}

		return config, nil
	}

	return RunConfig{}, fmt.Errorf("profile %s does not contain a definition for mode %s", p.name, mode)
}
