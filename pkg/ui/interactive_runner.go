package ui

import (
	"github.com/charmbracelet/bubbles/textinput"
	"github.com/zwoo-hq/zwooc/pkg/config"
)

type Model struct {
	inputOpen bool
	input     textinput.Model
}

// NewInteractiveRunner creates a new interactive runner for long running tasks
func NewInteractiveRunner(tasks config.TaskList, opts ViewOptions, conf config.Config) error {
	return nil
}
