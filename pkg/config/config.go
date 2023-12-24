package config

const (
	ModeRun   = "run"
	ModeBuild = "build"
	ModeWatch = "watch"
)

const (
	AdapterViteYarn = "vite-yarn"
	AdapterDotnet   = "dotnet"
)

const (
	KeyDefault  = "$default"
	KeyAdapter  = "$adapter"
	KeyFragment = "$fragments"
	KeyCompound = "$compounds"
	KeyPre      = "$pre"
	KeyPost     = "$post"
)

type (
	FragmentOptions map[string]string

	HookOptions struct {
		Fragments []string `json:"fragments"`
		Command   string   `json:"command"`
	}

	BaseOptions struct {
		Alias            string   `json:"alias"`
		SkipFragments    bool     `json:"skipFragments"`
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
		Profiles map[string]string `json:"profiles"`
	}
)

func IsReservedKey(key string) bool {
	switch key {
	case KeyAdapter:
		return true
	case KeyFragment:
		return true
	case KeyCompound:
		return true
	case KeyPre:
		return true
	case KeyPost:
		return true
	}
	return false
}

func IsValidRunMode(key string) bool {
	switch key {
	case ModeRun:
		return true
	case ModeBuild:
		return true
	case ModeWatch:
		return true
	}
	return false
}
