package config

import (
	"github.com/zwoo-hq/zwooc/pkg/helper"
	"github.com/zwoo-hq/zwooc/pkg/tasks"
)

type ResolvedFragment struct {
	Name       string
	Directory  string
	Command    string
	ProfileKey string
	Mode       string
	Options    map[string]interface{}
}

var _ Hookable = (*ResolvedFragment)(nil)

func (r ResolvedFragment) ResolvePreHook() ResolvedHook {
	if options, ok := r.Options[KeyPre]; ok {
		options := helper.MapToStruct(options.(map[string]interface{}), HookOptions{})
		return options.ResolveWithFragment(r, KeyPre)
	}
	return ResolvedHook{}
}

func (r ResolvedFragment) ResolvePostHook() ResolvedHook {
	if options, ok := r.Options[KeyPost]; ok {
		options := helper.MapToStruct(options.(map[string]interface{}), HookOptions{})
		return options.ResolveWithFragment(r, KeyPost)
	}
	return ResolvedHook{}
}

func (r ResolvedFragment) GetTask(extraArgs []string) tasks.Task {
	if r.Command == "" {
		return tasks.Empty()
	}
	return tasks.NewBasicCommandTask(r.Name, r.Command, r.Directory, extraArgs)
}

func (r ResolvedFragment) GetTaskWithBaseName(baseName string, extraArgs []string) tasks.Task {
	if r.Command == "" {
		return tasks.Empty()
	}
	return tasks.NewBasicCommandTask(helper.BuildName(baseName, r.Name), r.Command, r.Directory, extraArgs)
}
