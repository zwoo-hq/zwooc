package legacyui

import "github.com/zwoo-hq/zwooc/pkg/tasks"

func NewRunner(tasks tasks.TaskList, options ViewOptions) {
	if options.QuiteMode {
		newQuiteRunner(tasks, options)
		return
	}

	if options.DisableTUI {
		newStaticRunner(tasks, options)
		return
	}

	// try interactive view
	if err := newInteractiveRunner(tasks, options); err != nil {
		// fall back to static view
		newStaticRunner(tasks, options)
	}
}
