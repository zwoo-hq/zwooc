package config

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
