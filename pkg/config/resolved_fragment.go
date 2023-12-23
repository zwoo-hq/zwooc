package config

import (
	"github.com/zwoo-hq/zwooc/pkg/tasks"
)

type ResolvedFragment struct {
	Name       string
	Directory  string
	Command    string
	ProfileKey string
	Mode       string
	Options    map[string]interface{}
}

// func (r ResolvedFragment) GetPreHooks() HookOptions {
// 	if options, ok := r.Options[KeyPre]; ok {
// 		return helper.MapToStruct(options.(map[string]interface{}), HookOptions{})
// 	}
// 	return HookOptions{}
// }

// func (r ResolvedFragment) GetPostHooks() HookOptions {
// 	if options, ok := r.Options[KeyPost]; ok {
// 		return helper.MapToStruct(options.(map[string]interface{}), HookOptions{})
// 	}
// 	return HookOptions{}
// }

func (r ResolvedFragment) GetTask() (tasks.Task, error) {
	return tasks.NewBasicCommandTask(r.Name, r.Command, r.Directory), nil
}
