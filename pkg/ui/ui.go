package ui

import "github.com/zwoo-hq/zwooc/pkg/tasks"

func NewView(forest tasks.Collection, provider *SimpleStatusProvider, options ViewOptions) {
	if options.QuiteMode {
		newQuiteTreeView(forest, provider)
		return
	}

	if options.DisableTUI {
		newStaticTreeView(forest, provider, options)
		return
	}

	// try interactive view
	if err := newTreeProgressView(forest, provider, options); err != nil {
		// fall back to static view
		newStaticTreeView(forest, provider, options)
	}
}

func NewInteractiveView(forest tasks.Collection, provider *SchedulerStatusProvider, options ViewOptions) {
	if err := newInteractiveView(forest, provider, options); err != nil {
		// fall back to static view
		newStaticTreeView(forest, provider.SimpleStatusProvider, options)
	}
}
