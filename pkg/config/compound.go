package config

import (
	"github.com/zwoo-hq/zwooc/pkg/helper"
	"github.com/zwoo-hq/zwooc/pkg/model"
)

type Compound struct {
	name      string
	directory string
	raw       map[string]interface{}
}

func (c Compound) Name() string {
	return c.name
}

func (c Compound) ResolveConfig() ResolvedCompound {
	options := helper.MapToStruct(c.raw, model.CompoundOptions{})
	return ResolvedCompound{
		Name:             c.name,
		Directory:        c.directory,
		Profiles:         options.Profiles,
		IncludeFragments: options.IncludeFragments,
		Options:          c.raw,
	}
}
