package vite

import (
	"os"
	"os/exec"
	"strings"

	"github.com/zwoo-hq/zwooc/pkg/model"
	"github.com/zwoo-hq/zwooc/pkg/tasks"
)

func createViteTask(packageManager string, c model.ProfileWrapper, extraArgs []string) tasks.Task {
	cmd := exec.Command(packageManager)
	cmd.Args = append(cmd.Args, "vite")
	cmd.Args = append(cmd.Args, convertModeToVite(c.GetMode()))

	profileOptions := c.GetProfileOptions()
	cmd.Env = os.Environ()
	if os.Getenv("CI") != "true" {
		cmd.Env = append(cmd.Env, "FORCE_COLOR=1")
	} else {
		cmd.Env = append(cmd.Env, "NO_COLOR=true")
	}
	cmd.Env = append(cmd.Env, profileOptions.Env...)

	viteOptions := c.GetViteOptions()
	if viteOptions.Mode != "" {
		cmd.Args = append(cmd.Args, "--mode", viteOptions.Mode)
	}

	for k, v := range profileOptions.Args {
		if strings.HasPrefix(k, "-") {
			cmd.Args = append(cmd.Args, k, v)
		} else {
			cmd.Args = append(cmd.Args, "--"+k, v)
		}
	}

	cmd.Args = append(cmd.Args, extraArgs...)
	cmd.Dir = c.GetDirectory()
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
