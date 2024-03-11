package vite

import (
	"os"

	"github.com/zwoo-hq/zwooc/pkg/adapter/shared"
	"github.com/zwoo-hq/zwooc/pkg/model"
	"github.com/zwoo-hq/zwooc/pkg/tasks"
)

func createViteTask(packageManager string, c model.ProfileWrapper, extraArgs []string) tasks.Task {
	cmd, additionalArgs := shared.CreateBaseCommand(packageManager, c, extraArgs)
	cmd.Args = append(cmd.Args, "vite")
	cmd.Args = append(cmd.Args, convertModeToVite(c.GetMode()))

	if os.Getenv("CI") != "true" {
		cmd.Env = append(cmd.Env, "FORCE_COLOR=1")
	} else {
		cmd.Env = append(cmd.Env, "NO_COLOR=true")
	}

	viteOptions := c.GetViteOptions()
	if viteOptions.Mode != "" {
		cmd.Args = append(cmd.Args, "--mode", viteOptions.Mode)
	}

	cmd.Args = append(cmd.Args, additionalArgs...)
	return tasks.NewCommandTask(c.GetName(), cmd)
}

func convertModeToVite(mode string) string {
	switch mode {
	case model.ModeBuild:
		return "build"
	case model.ModeWatch:
		return "dev"
	case model.ModeRun:
		return "preview"
	}
	return ""
}
