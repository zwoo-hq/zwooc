package config

import "fmt"

type Fragment struct {
	name      string
	directory string
	raw       interface{}
}

func (f Fragment) Name() string {
	return f.name
}

func (f Fragment) GetConfig(mode string, callingProfile string) (ResolvedFragment, error) {
	if !IsValidRunMode(mode) && mode != "" {
		return ResolvedFragment{}, fmt.Errorf("invalid run mode: '%s'", mode)
	}

	if defaultCmd, ok := f.raw.(string); ok {
		return ResolvedFragment{
			Name:       f.name,
			Directory:  f.directory,
			Command:    defaultCmd,
			Options:    map[string]interface{}{},
			Mode:       mode,
			ProfileKey: callingProfile,
		}, nil
	}

	precedenceIndexes := []string{
		fmt.Sprintf("%s:%s", mode, callingProfile),
		callingProfile,
		mode,
		KeyDefault,
	}

	if options, ok := f.raw.(map[string]interface{}); ok {
		for _, index := range precedenceIndexes {
			if fragmentCommand, ok := options[index]; ok {
				return ResolvedFragment{
					Name:       f.name,
					Directory:  f.directory,
					Command:    fragmentCommand.(string),
					Options:    options,
					Mode:       mode,
					ProfileKey: callingProfile,
				}, nil
			}
		}
	}

	return ResolvedFragment{}, fmt.Errorf("fragment '%s' does not contain a definition for mode '%s'", f.name, mode)
}
