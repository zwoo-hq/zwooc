package config

import (
	"github.com/zwoo-hq/zwooc/pkg/helper"
	"github.com/zwoo-hq/zwooc/pkg/tasks"
)

type ResolvedHook struct {
	Kind      string
	Command   string
	Fragments []string
	Base      string
	Directory string
}

func (h HookOptions) ResolveWithProfile(callingProfile ResolvedProfile, kind string) ResolvedHook {
	return ResolvedHook{
		Kind:      kind,
		Command:   h.Command,
		Fragments: h.Fragments,
		Base:      helper.BuildName(callingProfile.Name, callingProfile.Mode),
		Directory: callingProfile.Directory,
	}
}

func (h HookOptions) ResolveWithFragment(callingFragment ResolvedFragment, kind string) ResolvedHook {
	return ResolvedHook{
		Kind:      kind,
		Command:   h.Command,
		Fragments: h.Fragments,
		Base:      callingFragment.Name,
		Directory: callingFragment.Directory,
	}
}

func (r ResolvedHook) GetTask() tasks.Task {
	if r.Command == "" {
		return tasks.Empty()
	}
	return tasks.NewBasicCommandTask(helper.BuildName(r.Base, r.Kind), r.Command, r.Directory, []string{})
}
