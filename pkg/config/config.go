package config

type Config struct {
	baseDir   string
	raw       map[string]interface{}
	profiles  []Profile
	fragments []Fragment
	compounds []Compound
}

func New(dir string, content map[string]interface{}) (Config, error) {
	c := Config{
		baseDir: dir,
		raw:     content,
	}
	err := c.init()
	return c, err
}

func NewLoaded(profiles []Profile, fragments []Fragment, compounds []Compound) Config {
	return Config{
		profiles:  profiles,
		fragments: fragments,
		compounds: compounds,
	}
}

const (
	ModeRun   = "run"
	ModeBuild = "build"
	ModeWatch = "watch"
)

const (
	AdapterViteYarn  = "vite-yarn"
	AdapterTauriYarn = "tauri-yarn"
	AdapterDotnet    = "dotnet"
)

const (
	KeyDefault   = "$default"
	KeyAdapter   = "$adapter"
	KeyDirectory = "$dir"
	KeyFragment  = "$fragments"
	KeyCompound  = "$compounds"
	KeyPre       = "$pre"
	KeyPost      = "$post"
)

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
	Hookable interface {
		ResolvePreHook() ResolvedHook
		ResolvePostHook() ResolvedHook
	}
)

func IsReservedKey(key string) bool {
	switch key {
	case KeyDefault:
		return true
	case KeyAdapter:
		return true
	case KeyDirectory:
		return true
	case KeyFragment:
		return true
	case KeyCompound:
		return true
	case KeyPre:
		return true
	case KeyPost:
		return true
	case "$schema":
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
