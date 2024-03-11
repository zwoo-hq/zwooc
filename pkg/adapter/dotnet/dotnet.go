package dotnet

import (
	"github.com/zwoo-hq/zwooc/pkg/adapter/shared"
	"github.com/zwoo-hq/zwooc/pkg/model"
	"github.com/zwoo-hq/zwooc/pkg/tasks"
)

type dotnetAdapter struct{}

var _ model.Adapter = (*dotnetAdapter)(nil)

func NewCliAdapter() model.Adapter {
	return &dotnetAdapter{}
}

func (a *dotnetAdapter) CreateTask(c model.ProfileWrapper, extraArgs []string) tasks.Task {
	cmd, additionalArgs := shared.CreateBaseCommand("dotnet", c, extraArgs)
	cmd.Args = append(cmd.Args, convertModeToDotnet(c.GetMode()))

	if c.GetMode() == model.ModeBuild {
		// run build mode by default in release mode
		cmd.Args = append(cmd.Args, "-c", "Release")
	}

	dotnetOptions := c.GetDotNetOptions()
	if dotnetOptions.Project != "" && c.GetMode() != model.ModeBuild {
		cmd.Args = append(cmd.Args, "--project", dotnetOptions.Project)
	}

	if dotnetOptions.Project != "" && c.GetMode() == model.ModeBuild {
		cmd.Args = append(cmd.Args, dotnetOptions.Project)
	}

	cmd.Args = append(cmd.Args, additionalArgs...)
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
