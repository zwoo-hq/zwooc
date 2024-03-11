package dotnet

import (
	"os"
	"os/exec"
	"strings"

	"github.com/zwoo-hq/zwooc/pkg/model"
	"github.com/zwoo-hq/zwooc/pkg/tasks"
)

type dotnetAdapter struct{}

var _ model.Adapter = (*dotnetAdapter)(nil)

func NewCliAdapter() model.Adapter {
	return &dotnetAdapter{}
}

func (a *dotnetAdapter) CreateTask(c model.ProfileWrapper, extraArgs []string) tasks.Task {
	cmd := exec.Command("dotnet")
	cmd.Args = append(cmd.Args, convertModeToDotnet(c.GetMode()))

	profileOptions := c.GetProfileOptions()
	cmd.Env = os.Environ()
	cmd.Env = append(cmd.Env, profileOptions.Env...)

	if c.GetMode() == model.ModeBuild {
		// run build mode by default in release mode
		cmd.Args = append(cmd.Args, "-c", "Release")
	}

	dotnetOptions := c.GetDotNetOptions()
	if dotnetOptions.Project != "" && c.GetMode() != model.ModeBuild {
		cmd.Args = append(cmd.Args, "--project", dotnetOptions.Project)
	}

	for k, v := range profileOptions.Args {
		if strings.HasPrefix(k, "-") {
			cmd.Args = append(cmd.Args, k, v)
		} else {
			cmd.Args = append(cmd.Args, "--"+k, v)
		}
	}

	cmd.Dir = c.GetDirectory()
	if dotnetOptions.Project != "" && c.GetMode() == model.ModeBuild {
		cmd.Args = append(cmd.Args, dotnetOptions.Project)
	}
	cmd.Args = append(cmd.Args, extraArgs...)
	return tasks.NewCommandTask(c.GetName(), cmd)
}

func convertModeToDotnet(mode string) string {
	switch mode {
	case model.ModeBuild:
		return "build"
	case model.ModeWatch:
		return "watch"
	case model.ModeRun:
		return "run"
	}
	return ""
}
