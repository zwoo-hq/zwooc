package config

import "fmt"

type Profile struct {
	name      string
	adapter   string
	directory string
	raw       map[string]interface{}
}

func (p Profile) Name() string {
	return p.name
}

func (p Profile) ResolveConfig(mode string) (ResolvedProfile, error) {
	if !IsValidRunMode(mode) {
		return ResolvedProfile{}, fmt.Errorf("invalid run mode: '%s'", mode)
	}

	options := p.raw[mode]
	if options == false {
		return ResolvedProfile{}, fmt.Errorf("profile '%s' disabled mode '%s'", p.name, mode)
	}

	config := ResolvedProfile{
		Name:      p.name,
		Adapter:   p.adapter,
		Directory: p.directory,
		Mode:      mode,
	}

	if optionsMap, ok := options.(map[string]interface{}); ok {
		config.Options = optionsMap
	}

	return config, nil
}
