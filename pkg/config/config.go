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
