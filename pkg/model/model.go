package model

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
