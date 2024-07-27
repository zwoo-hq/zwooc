package config

import (
	"slices"
)

type (
	LoadOptions struct {
		SkipHooks bool
		Exclude   []string
		ExtraArgs []string
	}

	RunnerOptions struct {
		MaxConcurrency  int
		UseLegacyRunner bool
	}

	loadingContext struct {
		skipHooks    bool
		excludedKeys []string
		extraArgs    []string
		callStack    []string
	}
)

func NewContext(opts LoadOptions) loadingContext {
	ctx := loadingContext{
		skipHooks:    opts.SkipHooks,
		excludedKeys: opts.Exclude,
		extraArgs:    opts.ExtraArgs,
		callStack:    []string{},
	}

	if ctx.excludedKeys == nil {
		ctx.excludedKeys = []string{}
	}
	if ctx.extraArgs == nil {
		ctx.extraArgs = []string{}
	}

	return ctx
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

func (c loadingContext) excludes(target string) bool {
	return slices.ContainsFunc(c.excludedKeys, func(s string) bool {
		return s == target
	})
}
