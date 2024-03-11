package config

import (
	"os"
	"os/exec"
	"strings"

	"github.com/zwoo-hq/zwooc/pkg/model"
	"github.com/zwoo-hq/zwooc/pkg/tasks"
)

func CreateTauriTask(c ResolvedProfile, extraArgs []string) tasks.Task {
	cmd := exec.Command("yarn")
	cmd.Args = append(cmd.Args, "tauri")
	cmd.Args = append(cmd.Args, convertModeToTauri(c.Mode))

	profileOptions := c.GetProfileOptions()
	cmd.Env = os.Environ()
	cmd.Env = append(cmd.Env, profileOptions.Env...)

	for k, v := range profileOptions.Args {
		if strings.HasPrefix(k, "-") {
			cmd.Args = append(cmd.Args, k, v)
		} else {
			cmd.Args = append(cmd.Args, "--"+k, v)
		}
	}

	cmd.Args = append(cmd.Args, extraArgs...)
	cmd.Dir = c.Directory
	return tasks.NewCommandTask(c.Name, cmd)
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
