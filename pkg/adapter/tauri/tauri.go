package tauri

import (
	"github.com/zwoo-hq/zwooc/pkg/adapter/shared"
	"github.com/zwoo-hq/zwooc/pkg/model"
	"github.com/zwoo-hq/zwooc/pkg/tasks"
)

func createTauriTask(packageManager string, c model.ProfileWrapper, extraArgs []string) tasks.Task {
	cmd, additionalArgs := shared.CreateBaseCommand(packageManager, c, extraArgs)
	cmd.Args = append(cmd.Args, "tauri")
	cmd.Args = append(cmd.Args, convertModeToTauri(c.GetMode()))

	cmd.Args = append(cmd.Args, additionalArgs...)
	return tasks.NewCommandTask(c.GetName(), cmd)
}

func convertModeToTauri(mode string) string {
	switch mode {
	case model.ModeBuild:
		return "build"
	case model.ModeWatch:
		return "dev"
	case model.ModeRun:
		return "dev"
	}
	return ""
}
