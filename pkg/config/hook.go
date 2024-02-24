package config

import (
	"github.com/zwoo-hq/zwooc/pkg/helper"
)

func (h HookOptions) ResolveWithProfile(callingProfile ResolvedProfile, kind string) ResolvedHook {
	return ResolvedHook{
		Kind:      kind,
		Command:   h.Command,
		Fragments: h.Fragments,
		Profiles:  h.Profiles,
		Base:      helper.BuildName(callingProfile.Name, callingProfile.Mode),
		Directory: callingProfile.Directory,
	}
}

func (h HookOptions) ResolveWithFragment(callingFragment ResolvedFragment, kind string) ResolvedHook {
	return ResolvedHook{
		Kind:      kind,
		Command:   h.Command,
		Fragments: h.Fragments,
		Profiles:  h.Profiles,
		Base:      callingFragment.Name,
		Directory: callingFragment.Directory,
	}
}

func (h HookOptions) ResolveWithCompound(callingCompound ResolvedCompound, kind string) ResolvedHook {
	return ResolvedHook{
		Kind:      kind,
		Command:   h.Command,
		Fragments: h.Fragments,
		Profiles:  h.Profiles,
		Base:      callingCompound.Name,
		Directory: callingCompound.Directory,
	}
}
