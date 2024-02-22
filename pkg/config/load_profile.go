package config

import (
	"fmt"

	"github.com/zwoo-hq/zwooc/pkg/helper"
	"github.com/zwoo-hq/zwooc/pkg/tasks"
	"golang.org/x/exp/maps"
)

func (c Config) ResolveProfile(key, mode string, extraArgs []string) (*tasks.TaskTreeNode, error) {
	if key == "" {
		key = KeyDefault
	}

	config, err := c.resolveRunConfig(key, mode)
	if err != nil {
		return nil, err
	}
	opts := config.GetBaseOptions()
	for opts.Base != "" {
		// load aliased profile
		newProfile, err := c.resolveRunConfig(opts.Base, mode)
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
	preStage, err := c.resolveHook(config.GetPreHooks(), config, config.Mode, config.Name)
	if err != nil {
		return nil, err
	}

	postStage, err := c.resolveHook(config.GetPostHooks(), config, config.Mode, config.Name)
	if err != nil {
		return nil, err
	}

	mainTask, err := config.GetTask(extraArgs)
	if err != nil {
		return nil, err
	}
	list := tasks.NewTaskTree(name, mainTask, mode == ModeWatch || mode == ModeRun)

	list.AddPreChild(preStage...)
	list.AddPostChild(postStage...)
	return list, nil
}

func (c Config) resolveRunConfig(key, mode string) (ResolvedProfile, error) {
	target, found := helper.FindBy(c.profiles, func(p Profile) bool {
		return p.Name() == key
	})
	if !found {
		return ResolvedProfile{}, fmt.Errorf("profile '%s' not found", key)
	}

	config, err := target.GetConfig(mode)
	if err != nil {
		return ResolvedProfile{}, err
	}
	return config, nil
}
