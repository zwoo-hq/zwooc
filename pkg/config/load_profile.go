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
	preStage, err := c.resolveHook(KeyPre, config, config.GetPreHooks())
	if err != nil {
		return TaskList{}, err
	}

	postStage, err := c.resolveHook(KeyPost, config, config.GetPostHooks())
	if err != nil {
		return TaskList{}, err
	}

	mainTask, err := config.GetTask()
	if err != nil {
		return TaskList{}, err
	}

	steps := []ExecutionStep{}
	if len(preStage) > 0 {
		steps = append(steps, ExecutionStep{
			Tasks:       preStage,
			Name:        helper.BuildName(name, KeyPre),
			RunParallel: true,
		})
	}
	steps = append(steps, ExecutionStep{
		Tasks: []tasks.Task{mainTask},
		Name:  name,
	})
	if len(postStage) > 0 {
		steps = append(steps, ExecutionStep{
			Tasks:       postStage,
			Name:        helper.BuildName(name, KeyPost),
			RunParallel: true,
		})
	}

	return TaskList{
		Name:  name,
		Steps: steps,
	}, nil
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

func (c Config) resolveHook(hookType string, profile ResolvedProfile, hook HookOptions) ([]tasks.Task, error) {
	baseName := hookType
	taskList := []tasks.Task{}
	if hook.Command != "" {
		taskList = append(taskList, tasks.NewBasicCommandTask(baseName, hook.Command, profile.Directory, []string{}))
	}

	for _, fragment := range hook.Fragments {
		fragmentConfig, err := c.resolveFragment(fragment, profile.Mode, profile.Name)
		if err != nil {
			return []tasks.Task{}, err
		}
		taskList = append(taskList, tasks.NewBasicCommandTask(helper.BuildName(baseName, fragment), fragmentConfig.Command, fragmentConfig.Directory, []string{}))
	}
	return taskList, nil
}
