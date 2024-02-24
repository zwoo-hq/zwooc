package config

import "github.com/zwoo-hq/zwooc/pkg/helper"

type ResolvedCompound struct {
	Name      string
	Directory string
	Profiles  map[string]string
	Options   map[string]interface{}
}

var _ Hookable = (*ResolvedCompound)(nil)

func (c ResolvedCompound) ResolvePreHook() ResolvedHook {
	if options, ok := c.Options[KeyPre]; ok {
		options := helper.MapToStruct(options.(map[string]interface{}), HookOptions{})
		return options.ResolveWithCompound(c, KeyPre)
	}
	return ResolvedHook{}
}

func (c ResolvedCompound) ResolvePostHook() ResolvedHook {
	if options, ok := c.Options[KeyPost]; ok {
		options := helper.MapToStruct(options.(map[string]interface{}), HookOptions{})
		return options.ResolveWithCompound(c, KeyPost)
	}
	return ResolvedHook{}
}
