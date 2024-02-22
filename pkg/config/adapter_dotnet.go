package config

import (
	"os"
	"os/exec"
	"strings"

	"github.com/zwoo-hq/zwooc/pkg/tasks"
)

func CreateDotnetTask(c ResolvedProfile, extraArgs []string) tasks.Task {
	cmd := exec.Command("dotnet")
	cmd.Args = append(cmd.Args, convertModeToDotnet(c.Mode))

	profileOptions := c.GetProfileOptions()
	cmd.Env = os.Environ()
	cmd.Env = append(cmd.Env, profileOptions.Env...)

	if c.Mode == ModeBuild {
		// run build mode by default in release mode
		cmd.Args = append(cmd.Args, "-c", "Release")
	}

	dotnetOptions := c.GetDotNetOptions()
	if dotnetOptions.Project != "" {
		cmd.Env = append(cmd.Args, "--project", dotnetOptions.Project)
	}

	for k, v := range profileOptions.Args {
		if strings.HasPrefix(k, "-") {
			cmd.Args = append(cmd.Args, k, v)
		} else {
			cmd.Args = append(cmd.Args, "--"+k, v)
		}
	}

	cmd.Dir = c.Directory
	cmd.Args = append(cmd.Args, extraArgs...)
	return tasks.NewCommandTask(c.Name, cmd)
}

func convertModeToDotnet(mode string) string {
	switch mode {
	case ModeBuild:
		return "build"
	case ModeWatch:
		return "watch"
	case ModeRun:
		return "run"
	}
	return ""
}
