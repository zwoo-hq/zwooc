package vite

import (
	"github.com/zwoo-hq/zwooc/pkg/model"
	"github.com/zwoo-hq/zwooc/pkg/tasks"
)

type viteAdapter struct {
	packageManager string
}

var _ model.Adapter = (*viteAdapter)(nil)

func NewYarnAdapter() model.Adapter {
	return &viteAdapter{"yarn"}
}

func NewNpmAdapter() model.Adapter {
	return &viteAdapter{"npm"}
}

func NewPnpmAdapter() model.Adapter {
	return &viteAdapter{"pnpm"}
}

func (a *viteAdapter) CreateTask(c model.ProfileWrapper, extraArgs []string) tasks.Task {
	return createViteTask(a.packageManager, c, extraArgs)
}
