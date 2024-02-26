package config

import (
	"fmt"

	"github.com/zwoo-hq/zwooc/pkg/helper"
	"github.com/zwoo-hq/zwooc/pkg/tasks"
)

func (c Config) loadAllHooks(caller Hookable, node *tasks.TaskTreeNode, mode, profile string, ctx loadingContext) error {
	preStage, err := c.loadHook(caller.ResolvePreHook(), mode, profile, ctx)
	if err != nil {
		return err
	}

	postStage, err := c.loadHook(caller.ResolvePostHook(), mode, profile, ctx)
	if err != nil {
		return err
	}

	node.AddPreChild(preStage...)
	node.AddPostChild(postStage...)
	return nil
}

var depth = 0

func (c Config) loadHook(hook ResolvedHook, mode, profile string, ctx loadingContext) ([]*tasks.TaskTreeNode, error) {
	depth += 1
	if depth > 1000 {
		return []*tasks.TaskTreeNode{}, fmt.Errorf("maximum depth of 1000 hooks reached (safety check for circular dependencies)")
	}

	ctx = ctx.withCaller(hook.Kind)
	taskList := []*tasks.TaskTreeNode{
		tasks.NewTaskTree(helper.BuildName(hook.Base, hook.Kind), hook.GetTask(), false),
	}

	for _, fragment := range hook.Fragments {
		if ctx.hasCaller(fragment) {
			return []*tasks.TaskTreeNode{}, createCircularDependencyError(ctx.callStack, fragment)
		}
		fragmentConfig, err := c.LoadFragment(combineFragmentKey(fragment, mode, profile), ctx)
		if err != nil {
			return nil, err
		}
		taskList = append(taskList, fragmentConfig)
	}

	for profile, mode := range hook.Profiles {
		if ctx.hasCaller(helper.BuildName(profile, mode)) {
			return []*tasks.TaskTreeNode{}, createCircularDependencyError(ctx.callStack, helper.BuildName(profile, mode))
		}
		profileConfig, err := c.LoadProfile(profile, mode, ctx)
		if err != nil {
			return nil, err
		}
		taskList = append(taskList, profileConfig)
	}

	return taskList, nil
}
