package tauri

import (
	"github.com/zwoo-hq/zwooc/pkg/model"
	"github.com/zwoo-hq/zwooc/pkg/tasks"
)

type tauriAdapter struct {
	packageManager string
}

var _ model.Adapter = (*tauriAdapter)(nil)

func NewYarnAdapter() model.Adapter {
	return &tauriAdapter{"yarn"}
}

func NewNpmAdapter() model.Adapter {
	return &tauriAdapter{"npm"}
}

func NewPnpmAdapter() model.Adapter {
	return &tauriAdapter{"pnpm"}
}

func (a *tauriAdapter) CreateTask(c model.ProfileWrapper, extraArgs []string) tasks.Task {
	return createTauriTask(a.packageManager, c, extraArgs)
}
