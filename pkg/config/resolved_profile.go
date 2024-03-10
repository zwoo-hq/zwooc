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
	if options, ok := r.Options[KeyPre]; ok {
		options := helper.MapToStruct(options.(map[string]interface{}), model.HookOptions{})
		return options.ResolveWithProfile(r, KeyPre)
	}
	return ResolvedHook{}
}

func (r ResolvedProfile) ResolvePostHook() ResolvedHook {
	if options, ok := r.Options[KeyPost]; ok {
		options := helper.MapToStruct(options.(map[string]interface{}), model.HookOptions{})
		return options.ResolveWithProfile(r, KeyPost)
	}
	return ResolvedHook{}
}

func (r ResolvedProfile) GetTask(args []string) (tasks.Task, error) {
	switch r.Adapter {
	case AdapterViteYarn:
		return CreateViteTask(r, args), nil
	case AdapterDotnet:
		return CreateDotnetTask(r, args), nil
	case AdapterTauriYarn:
		return CreateTauriTask(r, args), nil
	}
	return tasks.Empty(), fmt.Errorf("unknown adapter: '%s'", r.Adapter)
}
