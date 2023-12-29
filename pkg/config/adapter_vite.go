package config

import (
	"os"
	"os/exec"

	"github.com/zwoo-hq/zwooc/pkg/tasks"
)

func CreateViteTask(c ResolvedProfile) tasks.Task {
	cmd := exec.Command("yarn")
	cmd.Args = append(cmd.Args, "vite")
	cmd.Args = append(cmd.Args, convertModeToVite(c.Mode))
	cmd.Stdout = os.Stdout

	profileOptions := c.GetProfileOptions()
	cmd.Env = os.Environ()
	cmd.Env = append(cmd.Env, "FORCE_COLOR=1")
	cmd.Env = append(cmd.Env, profileOptions.Env...)
	for k, v := range profileOptions.Args {
		cmd.Args = append(cmd.Args, "--"+k)
		cmd.Args = append(cmd.Args, v)
	}

	cmd.Dir = c.Directory
	return tasks.NewCommandTask(c.Name, cmd)
}

func convertModeToVite(mode string) string {
	switch mode {
	case ModeBuild:
		return "build"
	case ModeWatch:
		return "dev"
	case ModeRun:
		return "preview"
	}
	return ""
}
