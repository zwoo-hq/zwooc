package config

import (
	"fmt"

	"github.com/zwoo-hq/zwooc/pkg/helper"
	"github.com/zwoo-hq/zwooc/pkg/tasks"
	"golang.org/x/exp/maps"
)

func (c Config) ResolveProfile(key, mode string) (TaskList, error) {
	if key == "" {
		key = KeyDefault
	}

	config, err := c.resolveRunConfig(key, mode)
	if err != nil {
		return TaskList{}, err
	}
	opts := config.GetBaseOptions()
	for opts.Alias != "" {
		// load aliased profile
		newProfile, err := c.resolveRunConfig(opts.Alias, mode)
		if err != nil {
			return TaskList{}, err
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
		if newProfile.GetBaseOptions().Alias == "" {
			break
		}
	}

	name := helper.BuildName(key, mode)
	preStage, err := c.resolveHook(config.GetPreHooks(), config, config.Mode, config.Name)
	if err != nil {
		return TaskList{}, err
	}

	postStage, err := c.resolveHook(config.GetPostHooks(), config, config.Mode, config.Name)
	if err != nil {
		return TaskList{}, err
	}

	mainTask, err := config.GetTask()
	if err != nil {
		return TaskList{}, err
	}
	list := NewTaskList(name, []ExecutionStep{
		{
			Name:  name,
			Tasks: []tasks.Task{mainTask},
		},
	})

	list.InsertBefore(preStage)
	list.InsertAfter(postStage)
	list.RemoveEmptyStagesAndTasks()
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

func (c Config) resolveHook(hook ResolvedHook, caller Hookable, mode, profile string) (TaskList, error) {
	taskList := []tasks.Task{hook.GetTask()}
	for _, fragment := range hook.Fragments {
		fragmentConfig, err := c.resolveFragment(fragment, mode, profile)
		if err != nil {
			return TaskList{}, err
		}
		taskList = append(taskList, fragmentConfig.GetTaskWithBaseName(helper.BuildName(hook.Base, hook.Kind), []string{}))
	}

	return NewTaskList("", []ExecutionStep{
		{
			Name:        helper.BuildName(hook.Base, hook.Kind),
			Tasks:       taskList,
			RunParallel: true,
		},
	}), nil
}
