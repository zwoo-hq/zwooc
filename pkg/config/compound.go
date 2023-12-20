package config

type Compound struct {
	name string
	raw  map[string]interface{}
}

func newCompound(name string, data map[string]interface{}) Compound {
	c := Compound{
		name: name,
		raw:  data,
	}
	return c
}

func (c Compound) Name() string {
	return c.name
}
