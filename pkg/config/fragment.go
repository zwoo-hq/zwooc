package config

type Fragment struct {
	name string
	raw  map[string]interface{}
}

func newFragment(name string, data map[string]interface{}) Fragment {
	f := Fragment{
		name: name,
		raw:  data,
	}
	return f
}

func (f Fragment) Name() string {
	return f.name
}
