package ui

import "github.com/zwoo-hq/zwooc/pkg/tasks"

func NewRunner(forest tasks.Collection, options ViewOptions) {
	if options.QuiteMode {
		newQuiteRunner(forest, options)
		return
	}

	if options.DisableTUI {
		newStaticTreeRunner(forest, options)
		return
	}

	// try interactive view
	// if err := NewStatusView(task, options); err != nil {
	// 	// fall back to static view
	// 	newStaticRunner(task, options)
	// }
}
