package config

import (
	"github.com/zwoo-hq/zwooc/pkg/helper"
	"github.com/zwoo-hq/zwooc/pkg/tasks"
)

type ResolvedHook struct {
	Kind      string
	Command   string
	Fragments []string
	Profiles  map[string]string
	Base      string
	Directory string
}

func (r ResolvedHook) GetTask() tasks.Task {
	if r.Command == "" {
		return tasks.Empty()
	}
	return tasks.NewBasicCommandTask(helper.BuildName(r.Base, r.Kind), r.Command, r.Directory, []string{})
}
