package config

import (
	"github.com/zwoo-hq/zwooc/pkg/helper"
	"github.com/zwoo-hq/zwooc/pkg/model"
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
	if options, ok := r.Options[model.KeyPre]; ok {
		hook := Hook{options.(map[string]interface{})}
		return hook.ResolveWithFragment(r, model.KeyPre)
	}
	return ResolvedHook{}
}

func (r ResolvedFragment) ResolvePostHook() ResolvedHook {
	if options, ok := r.Options[model.KeyPost]; ok {
		hook := Hook{options.(map[string]interface{})}
		return hook.ResolveWithFragment(r, model.KeyPost)
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
