package config

import (
	"fmt"

	"github.com/zwoo-hq/zwooc/pkg/helper"
	"github.com/zwoo-hq/zwooc/pkg/tasks"
)

func (c Config) LoadCompound(key string, ctx loadingContext) (tasks.Collection, error) {
	if ctx.excludes(key) {
		return nil, ErrTargetExcluded
	}

	compound, err := c.resolveCompound(key)
	if err != nil {
		return []*tasks.TaskTreeNode{}, err
	}

	nodes := tasks.NewCollection()
	compoundNode := tasks.NewTaskTree(key, tasks.Empty(), false)
	if !ctx.skipHooks {
		// sue the compound key as profile here
		err = c.loadAllHooks(compound, compoundNode, "", key, ctx)
		if err != nil {
			return nil, err
		}
	}
	nodes = append(nodes, compoundNode)

	for profileKey, mode := range compound.Profiles {
		resolved, err := c.LoadProfile(profileKey, mode, ctx.withCaller(key))
		if err != nil {
			return []*tasks.TaskTreeNode{}, err
		}
		nodes = append(nodes, resolved...)
	}

	for _, fragmentKey := range compound.IncludeFragments {
		fragment, err := c.LoadFragment(combineFragmentKey(fragmentKey, "", key), ctx.withCaller("includes"))
		if err != nil {
			return nil, err
		}
		nodes = append(nodes, fragment)
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
