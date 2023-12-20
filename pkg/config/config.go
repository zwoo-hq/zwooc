package config

const (
	ModeRun   = "run"
	ModeBuild = "build"
	ModeWatch = "watch"
)

const (
	KeyDefault  = "$default"
	KeyAdapter  = "$adapter"
	KeyFragment = "$fragment"
	KeyCompound = "$compound"
	KeyPre      = "$pre"
	KeyPost     = "$post"
)

type (
	Fragment map[string]string

	HookOptions struct {
		Fragments []string `json:"fragments"`
		Command   string   `json:"command"`
	}

	BaseOptions struct {
		Alias            string   `json:"alias"`
		SkipFragments    bool     `json:"skipFragments"`
		IncludeFragments []string `json:"includeFragments"`
	}

	BaseProfile struct {
		Args map[string]string `json:"args"`
		Env  []string          `json:"env"`
	}

	ViteOptions struct {
		Mode string `json:"mode"`
	}

	DotNetOptions struct {
		Project string `json:"project"`
	}

	Compound struct {
		Profiles map[string]string `json:"profiles"`
	}
)

func IsReservedKey(key string) bool {
	switch key {
	case KeyDefault:
	case KeyAdapter:
	case KeyFragment:
	case KeyCompound:
	case KeyPre:
	case KeyPost:
		return true
	}
	return false
}
