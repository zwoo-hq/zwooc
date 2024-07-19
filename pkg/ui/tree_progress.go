package ui

import (
	"github.com/zwoo-hq/zwooc/pkg/tasks"
)

type TreeStatusProvider interface {
	Status() []*TreeStatusNode
}

type TreeProgressView struct {
	status  []*TreeStatusNode
	opts    ViewOptions
	outputs map[string]*tasks.CommandCapturer
}

type TreeStatusNode struct {
	ID       string
	Name     string
	Error    error
	Status   TaskStatus
	Children []*TreeStatusNode
	Task     tasks.Task
}

func NewTreeProgressView(status TreeStatusProvider, opts ViewOptions) *TreeProgressView {
	return &TreeProgressView{
		status: status.Status(),
		opts:   opts,
	}
}
