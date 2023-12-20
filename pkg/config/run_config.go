package config

import "github.com/zwoo-hq/zwooc/pkg/helper"

type RunConfig struct {
	Name      string
	Adapter   string
	Directory string
	Options   map[string]interface{}
}

func (r RunConfig) GetViteOptions() ViteOptions {
	return helper.MapToStruct(r.Options, ViteOptions{})
}

func (r RunConfig) GetDotNetOptions() DotNetOptions {
	return helper.MapToStruct(r.Options, DotNetOptions{})
}

func (r RunConfig) GetBaseOptions() BaseOptions {
	return helper.MapToStruct(r.Options, BaseOptions{})
}

func (r RunConfig) GetProfileOptions() ProfileOptions {
	return helper.MapToStruct(r.Options, ProfileOptions{})
}

func (r RunConfig) GetPreHooks() HookOptions {
	if options, ok := r.Options[KeyPre]; ok {
		return helper.MapToStruct(options.(map[string]interface{}), HookOptions{})
	}
	return HookOptions{}
}

func (r RunConfig) GetPostHooks() HookOptions {
	if options, ok := r.Options[KeyPost]; ok {
		return helper.MapToStruct(options.(map[string]interface{}), HookOptions{})
	}
	return HookOptions{}
}
