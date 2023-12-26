package ui

import "github.com/zwoo-hq/zwooc/pkg/config"

func NewRunner(tasks config.TaskList, options ViewOptions) {
	newStaticRunner(tasks)
	return
	// try interactive view
	if err := newInteractiveRunner(tasks); err != nil {
		// fall back to static view
	}
}
