package config

import "github.com/zwoo-hq/zwooc/pkg/model"

type ResolvedCompound struct {
	Name             string
	Directory        string
	Profiles         map[string]string
	IncludeFragments []string
	Options          map[string]interface{}
}

var _ Hookable = (*ResolvedCompound)(nil)

func (c ResolvedCompound) ResolvePreHook() ResolvedHook {
	if options, ok := c.Options[model.KeyPre]; ok {
		hook := Hook{options.(map[string]interface{})}
		return hook.ResolveWithCompound(c, model.KeyPre)
	}
	return ResolvedHook{}
}

func (c ResolvedCompound) ResolvePostHook() ResolvedHook {
	if options, ok := c.Options[model.KeyPost]; ok {
		hook := Hook{options.(map[string]interface{})}
		return hook.ResolveWithCompound(c, model.KeyPost)
	}
	return ResolvedHook{}
}
