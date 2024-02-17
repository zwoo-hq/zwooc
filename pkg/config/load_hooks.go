package config

import (
	"github.com/zwoo-hq/zwooc/pkg/helper"
	"github.com/zwoo-hq/zwooc/pkg/tasks"
)

func (c Config) resolveHooks(caller Hookable, node *tasks.TaskTreeNode, mode, profile string) error {
	preStage, err := c.resolveHook(caller.GetPreHooks(), caller, mode, profile)
	if err != nil {
		return err
	}

	postStage, err := c.resolveHook(caller.GetPostHooks(), caller, mode, profile)
	if err != nil {
		return err
	}

	node.AddPreChild(preStage...)
	node.AddPostChild(postStage...)
	return nil
}

func (c Config) resolveHook(hook ResolvedHook, caller Hookable, mode, profile string) ([]*tasks.TaskTreeNode, error) {
	taskList := []*tasks.TaskTreeNode{
		tasks.NewTaskTree(helper.BuildName(hook.Base, hook.Kind), hook.GetTask(), false),
	}

	for _, fragment := range hook.Fragments {
		fragmentConfig, err := c.resolveFragment(fragment, mode, profile)
		if err != nil {
			return nil, err
		}
		name := helper.BuildName(hook.Base, hook.Kind)
		taskList = append(taskList, tasks.NewTaskTree(fragmentConfig.Name, fragmentConfig.GetTaskWithBaseName(name, []string{}), false))
	}

	for profile, mode := range hook.Profiles {
		profileConfig, err := c.ResolveProfile(profile, mode)
		if err != nil {
			return nil, err
		}
		taskList = append(taskList, profileConfig)
	}

	return taskList, nil
}
