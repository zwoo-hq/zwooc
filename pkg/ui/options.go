package ui

import "github.com/zwoo-hq/zwooc/pkg/tasks"

type ViewOptions struct {
	DisableTUI    bool
	QuiteMode     bool
	InlineOutput  bool
	CombineOutput bool
	DisablePrefix bool
	// todo: move out of here
	MaxConcurrency int

	Forest         tasks.Collection
	StatusProvider GenericStatusProvider
}
