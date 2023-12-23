package config

import "fmt"

type Fragment struct {
	name      string
	directory string
	raw       map[string]interface{}
}

func newFragment(name, directory string, data map[string]interface{}) Fragment {
	f := Fragment{
		name:      name,
		directory: directory,
		raw:       data,
	}
	return f
}

func (f Fragment) Name() string {
	return f.name
}

func (f Fragment) GetConfig(mode string, callingProfile string) (ResolvedFragment, error) {
	if !IsValidRunMode(mode) {
		return ResolvedFragment{}, fmt.Errorf("invalid run mode: %s", mode)
	}

	precedenceIndexes := []string{
		fmt.Sprintf("%s:%s", mode, callingProfile),
		mode,
		KeyDefault,
	}

	for _, index := range precedenceIndexes {
		if fragmentCommand, ok := f.raw[index]; ok {
			return ResolvedFragment{
				Name:       f.name,
				Directory:  f.directory,
				Command:    fragmentCommand.(string),
				Options:    f.raw,
				Mode:       mode,
				ProfileKey: callingProfile,
			}, nil
		}
	}

	return ResolvedFragment{}, fmt.Errorf("fragment %s does not contain a definition for mode %s", f.name, mode)
}
