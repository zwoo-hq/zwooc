package config

import (
	"fmt"

	"github.com/zwoo-hq/zwooc/pkg/helper"
	"github.com/zwoo-hq/zwooc/pkg/model"
	"github.com/zwoo-hq/zwooc/pkg/tasks"
)

type ResolvedProfile struct {
	Name      string
	Mode      string
	Adapter   string
	Directory string
	Options   map[string]interface{}
}

var _ Hookable = (*ResolvedProfile)(nil)
var _ model.ProfileWrapper = (*ResolvedProfile)(nil)

func (r ResolvedProfile) GetName() string {
	return r.Name
}

func (r ResolvedProfile) GetMode() string {
	return r.Mode
}

func (r ResolvedProfile) GetAdapter() string {
	return r.Adapter
}

func (r ResolvedProfile) GetDirectory() string {
	return r.Directory
}

func (r ResolvedProfile) GetOptions() map[string]interface{} {
	return r.Options
}

func (r ResolvedProfile) GetViteOptions() model.ViteOptions {
	return helper.MapToStruct(r.Options, model.ViteOptions{})
}

func (r ResolvedProfile) GetDotNetOptions() model.DotNetOptions {
	return helper.MapToStruct(r.Options, model.DotNetOptions{})
}

func (r ResolvedProfile) GetBaseOptions() model.BaseOptions {
	return helper.MapToStruct(r.Options, model.BaseOptions{})
}

func (r ResolvedProfile) GetProfileOptions() model.ProfileOptions {
	return helper.MapToStruct(r.Options, model.ProfileOptions{})
}

func (r ResolvedProfile) ResolvePreHook() ResolvedHook {
	if options, ok := r.Options[model.KeyPre]; ok {
		hook := Hook{options.(map[string]interface{})}
		return hook.ResolveWithProfile(r, model.KeyPre)
	}
	return ResolvedHook{}
}

func (r ResolvedProfile) ResolvePostHook() ResolvedHook {
	if options, ok := r.Options[model.KeyPost]; ok {
		hook := Hook{options.(map[string]interface{})}
		return hook.ResolveWithProfile(r, model.KeyPost)
	}
	return ResolvedHook{}
}

func (r ResolvedProfile) GetTask(args []string) (tasks.Task, error) {
	adapter := GetAdapter(r.Adapter)
	if adapter != nil {
		return adapter.CreateTask(r, args), nil
	}
	return tasks.Empty(), fmt.Errorf("unknown adapter: '%s'", r.Adapter)
}
