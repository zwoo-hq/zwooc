package vite

import (
	"github.com/zwoo-hq/zwooc/pkg/model"
	"github.com/zwoo-hq/zwooc/pkg/tasks"
)

type viteAdapter struct{}

var _ model.Adapter = (*viteAdapter)(nil)

func NewYarnAdapter() model.Adapter {
	return &viteAdapter{}
}

func (a *viteAdapter) CreateTask(c model.ProfileWrapper, extraArgs []string) tasks.Task {
	return createViteTask("yarn", c, extraArgs)
}
