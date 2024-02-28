package config

import (
	"fmt"
	"strings"
)

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
		Profiles map[string]string `json:"profiles"`
	}
)

type (
	Hookable interface {
		ResolvePreHook() ResolvedHook
		ResolvePostHook() ResolvedHook
	}

	LoadOptions struct {
		SkipHooks bool
		Exclude   []string
		ExtraArgs []string
	}

	loadingContext struct {
		skipHooks bool
		exclude   []string
		extraArgs []string
		callStack []string
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

func NewContext(opts LoadOptions) loadingContext {
	return loadingContext{
		skipHooks: opts.SkipHooks,
		exclude:   opts.Exclude,
		extraArgs: opts.ExtraArgs,
		callStack: []string{},
	}
}

func (c loadingContext) getArgs() []string {
	if len(c.callStack) == 0 {
		return c.extraArgs
	}
	return []string{}
}

func (c loadingContext) withCaller(caller string) loadingContext {
	c.callStack = append(c.callStack, caller)
	return c
}

func (c loadingContext) hasCaller(caller string) bool {
	for _, c := range c.callStack {
		if c == caller {
			return true
		}
	}
	return false
}

func createCircularDependencyError(caller []string, target string) error {
	return fmt.Errorf("circular dependency detected: '%s' from %s", target, strings.Join(caller, " -> "))
}
