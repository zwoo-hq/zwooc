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

func (r ResolvedProfile) GetPreHooks() HookOptions {
	if options, ok := r.Options[KeyPre]; ok {
		return helper.MapToStruct(options.(map[string]interface{}), HookOptions{})
	}
	return HookOptions{}
}

func (r ResolvedProfile) GetPostHooks() HookOptions {
	if options, ok := r.Options[KeyPost]; ok {
		return helper.MapToStruct(options.(map[string]interface{}), HookOptions{})
	}
	return HookOptions{}
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
