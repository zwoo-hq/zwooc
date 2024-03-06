package config

type (
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
