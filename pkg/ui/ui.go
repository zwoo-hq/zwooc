package ui

import (
	"github.com/zwoo-hq/zwooc/pkg/config"
)

func NewRunner(tasks config.TaskList, options ViewOptions) {
	if options.QuiteMode {
		newQuiteRunner(tasks)
		return
	}

	if options.DisableTUI {
		newStaticRunner(tasks, options)
		return
	}

	// try interactive view
	if err := newInteractiveRunner(tasks); err != nil {
		// fall back to static view
		newStaticRunner(tasks, options)
	}
}
