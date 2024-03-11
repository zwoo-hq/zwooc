package model

import "github.com/zwoo-hq/zwooc/pkg/tasks"

type (
	FragmentOptions map[string]interface{}

	HookOptions struct {
		Command   string            `json:"command"`
		Fragments []string          `json:"fragments"`
		Profiles  map[string]string `json:"profiles"`
	}

	BaseOptions struct {
		Base             string   `json:"base"`
		IncludeFragments []string `json:"includeFragments"`
	}

	ProfileOptions struct {
		Args map[string]string `json:"args"`
		Env  []string          `json:"env"`
	}

	ViteOptions struct {
		Mode string `json:"mode"`
	}

	DotNetOptions struct {
		Project string `json:"project"`
	}

	CompoundOptions struct {
		Profiles         map[string]string `json:"profiles"`
		IncludeFragments []string          `json:"includeFragments"`
	}
)

type (
	ProfileWrapper interface {
		GetName() string
		GetMode() string
		GetAdapter() string
		GetDirectory() string
		GetOptions() map[string]interface{}
		GetViteOptions() ViteOptions
		GetDotNetOptions() DotNetOptions
		GetBaseOptions() BaseOptions
		GetProfileOptions() ProfileOptions
	}

	Adapter interface {
		CreateTask(c ProfileWrapper, extraArgs []string) tasks.Task
	}
)
