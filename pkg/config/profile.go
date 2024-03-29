package config

import (
	"fmt"

	"github.com/zwoo-hq/zwooc/pkg/helper"
	"github.com/zwoo-hq/zwooc/pkg/model"
)

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

	if command, ok := options.(string); ok && p.adapter == model.AdapterCustom {
		// support custom profiles
		options = map[string]interface{}{"command": command}
	}

	config := ResolvedProfile{
		Name:      p.name,
		Adapter:   p.adapter,
		Directory: p.directory,
		Mode:      mode,
		Options:   map[string]interface{}{},
	}

	if optionsMap, ok := options.(map[string]interface{}); ok {
		config.Options = optionsMap
	}

	// hoist "global" options
	var allowedHoistedOptions = map[string]interface{}{}
	for optionKey, optionValue := range p.raw {
		if !IsValidRunMode(optionKey) {
			allowedHoistedOptions[optionKey] = optionValue
		}
	}
	config.Options = helper.MergeDeep(allowedHoistedOptions, config.Options)

	return config, nil
}
