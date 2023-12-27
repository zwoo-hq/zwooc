package config

import (
	"fmt"

	"github.com/zwoo-hq/zwooc/pkg/helper"
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

func (r ResolvedProfile) GetViteOptions() ViteOptions {
	return helper.MapToStruct(r.Options, ViteOptions{})
}

func (r ResolvedProfile) GetDotNetOptions() DotNetOptions {
	return helper.MapToStruct(r.Options, DotNetOptions{})
}

func (r ResolvedProfile) GetBaseOptions() BaseOptions {
	return helper.MapToStruct(r.Options, BaseOptions{})
}

func (r ResolvedProfile) GetProfileOptions() ProfileOptions {
	return helper.MapToStruct(r.Options, ProfileOptions{})
}

func (r ResolvedProfile) GetPreHooks() ResolvedHook {
	if options, ok := r.Options[KeyPre]; ok {
		options := helper.MapToStruct(options.(map[string]interface{}), HookOptions{})
		return options.ResolveWithProfile(r, KeyPre)
	}
	return ResolvedHook{}
}

func (r ResolvedProfile) GetPostHooks() ResolvedHook {
	if options, ok := r.Options[KeyPost]; ok {
		options := helper.MapToStruct(options.(map[string]interface{}), HookOptions{})
		return options.ResolveWithProfile(r, KeyPost)
	}
	return ResolvedHook{}
}

func (r ResolvedProfile) GetTask() (tasks.Task, error) {
	switch r.Adapter {
	case AdapterViteYarn:
		return CreateViteTask(r), nil
	case AdapterDotnet:
		return tasks.Empty(), nil
	}
	return tasks.Empty(), fmt.Errorf("unknown adapter: '%s'", r.Adapter)
}
