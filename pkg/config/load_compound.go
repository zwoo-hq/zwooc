package config

import (
	"fmt"

	"github.com/zwoo-hq/zwooc/pkg/helper"
	"github.com/zwoo-hq/zwooc/pkg/tasks"
)

func (c Config) LoadCompound(key string) ([]*tasks.TaskTreeNode, error) {
	compound, err := c.resolveCompound(key)
	if err != nil {
		return []*tasks.TaskTreeNode{}, err
	}

	nodes := []*tasks.TaskTreeNode{}
	for profileKey, mode := range compound.Profiles {
		resolved, err := c.LoadProfile(profileKey, mode, []string{})
		if err != nil {
			return []*tasks.TaskTreeNode{}, err
		}
		nodes = append(nodes, resolved)
	}

	return nodes, nil
}

func (c Config) resolveCompound(key string) (ResolvedCompound, error) {
	target, found := helper.FindBy(c.compounds, func(c Compound) bool {
		return c.Name() == key
	})
	if !found {
		return ResolvedCompound{}, fmt.Errorf("compound '%s' not found", key)
	}

	return target.ResolveConfig(), nil
}
