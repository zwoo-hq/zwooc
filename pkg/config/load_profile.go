package config

import (
	"fmt"

	"github.com/zwoo-hq/zwooc/pkg/helper"
	"github.com/zwoo-hq/zwooc/pkg/tasks"
)

func (c Config) ResolveProfile(key, mode string) (TaskList, error) {
	config, err := c.resolveRunConfig(key, mode)
	if err != nil {
		return TaskList{}, err
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
	if key == "" {
		key = KeyDefault
	}

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
