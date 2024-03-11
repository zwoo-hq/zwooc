package config

import (
	"github.com/zwoo-hq/zwooc/pkg/helper"
	"github.com/zwoo-hq/zwooc/pkg/model"
)

type Hook struct {
	raw map[string]interface{}
}

func (h Hook) getOptions() model.HookOptions {
	return helper.MapToStruct(h.raw, model.HookOptions{})
}

func (h Hook) ResolveWithProfile(callingProfile ResolvedProfile, kind string) ResolvedHook {
	options := h.getOptions()

	return ResolvedHook{
		Kind:      kind,
		Command:   options.Command,
		Fragments: options.Fragments,
		Profiles:  options.Profiles,
		Base:      helper.BuildName(callingProfile.Name, callingProfile.Mode),
		Directory: callingProfile.Directory,
	}
}

func (h Hook) ResolveWithFragment(callingFragment ResolvedFragment, kind string) ResolvedHook {
	options := h.getOptions()

	return ResolvedHook{
		Kind:      kind,
		Command:   options.Command,
		Fragments: options.Fragments,
		Profiles:  options.Profiles,
		Base:      callingFragment.Name,
		Directory: callingFragment.Directory,
	}
}

func (h Hook) ResolveWithCompound(callingCompound ResolvedCompound, kind string) ResolvedHook {
	options := h.getOptions()

	return ResolvedHook{
		Kind:      kind,
		Command:   options.Command,
		Fragments: options.Fragments,
		Profiles:  options.Profiles,
		Base:      callingCompound.Name,
		Directory: callingCompound.Directory,
	}
}
