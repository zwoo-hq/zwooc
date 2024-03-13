package config

import (
	"github.com/zwoo-hq/zwooc/pkg/adapter/custom"
	"github.com/zwoo-hq/zwooc/pkg/adapter/dotnet"
	"github.com/zwoo-hq/zwooc/pkg/adapter/tauri"
	"github.com/zwoo-hq/zwooc/pkg/adapter/vite"
	"github.com/zwoo-hq/zwooc/pkg/model"
)

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

type (
	Hookable interface {
		ResolvePreHook() ResolvedHook
		ResolvePostHook() ResolvedHook
	}
)

func IsReservedKey(key string) bool {
	switch key {
	case model.KeyDefault:
		return true
	case model.KeyAdapter:
		return true
	case model.KeyDirectory:
		return true
	case model.KeyFragment:
		return true
	case model.KeyCompound:
		return true
	case model.KeyPre:
		return true
	case model.KeyPost:
		return true
	case "$schema":
		return true
	}
	return false
}

func IsValidRunMode(key string) bool {
	switch key {
	case model.ModeRun:
		return true
	case model.ModeBuild:
		return true
	case model.ModeWatch:
		return true
	}
	return false
}

func GetAdapter(adapter string) model.Adapter {
	switch adapter {
	case model.AdapterViteYarn:
		return vite.NewYarnAdapter()
	case model.AdapterViteNpm:
		return vite.NewNpmAdapter()
	case model.AdapterVitePnpm:
		return vite.NewPnpmAdapter()
	case model.AdapterTauriYarn:
		return tauri.NewYarnAdapter()
	case model.AdapterTauriNpm:
		return tauri.NewNpmAdapter()
	case model.AdapterTauriPnpm:
		return tauri.NewPnpmAdapter()
	case model.AdapterDotnet:
		return dotnet.NewCliAdapter()
	case model.AdapterCustom:
		return custom.NewAdapter()
	}
	return nil
}
