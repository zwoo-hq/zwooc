package config

import (
	"fmt"

	"github.com/zwoo-hq/zwooc/pkg/helper"
	"github.com/zwoo-hq/zwooc/pkg/tasks"
	"golang.org/x/exp/maps"
)

func (c Config) LoadProfile(key, mode string, ctx loadingContext) (*tasks.TaskTreeNode, error) {
	if key == "" {
		key = KeyDefault
	}

	config, err := c.resolveProfile(key, mode)
	if err != nil {
		return nil, err
	}
	opts := config.GetBaseOptions()
	for opts.Base != "" {
		// load aliased profile
		newProfile, err := c.resolveProfile(opts.Base, mode)
		if err != nil {
			return nil, err
		}
		// merge profiles
		config = ResolvedProfile{
			Name:      config.Name,
			Mode:      config.Mode,
			Adapter:   newProfile.Adapter,
			Directory: newProfile.Directory,
			Options:   helper.MergeDeep(maps.Clone(newProfile.Options), config.Options),
		}
		opts = config.GetBaseOptions()
		if newProfile.GetBaseOptions().Base == "" {
			break
		}
	}

	name := helper.BuildName(key, mode)
	mainTask, err := config.GetTask(ctx.getArgs())
	if err != nil {
		return nil, err
	}
	treeNode := tasks.NewTaskTree(name, mainTask, mode == ModeWatch || mode == ModeRun)
	err = c.loadAllHooks(config, treeNode, mode, key, ctx.withCaller(name))
	if err != nil {
		return nil, err
	}
	return treeNode, nil
}

func (c Config) resolveProfile(key, mode string) (ResolvedProfile, error) {
	target, found := helper.FindBy(c.profiles, func(p Profile) bool {
		return p.Name() == key
	})
	if !found {
		return ResolvedProfile{}, fmt.Errorf("profile '%s' not found", key)
	}

	config, err := target.ResolveConfig(mode)
	if err != nil {
		return ResolvedProfile{}, err
	}
	return config, nil
}
