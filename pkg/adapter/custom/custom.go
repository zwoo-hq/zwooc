package custom

import (
	"strings"

	"github.com/zwoo-hq/zwooc/pkg/adapter/shared"
	"github.com/zwoo-hq/zwooc/pkg/helper"
	"github.com/zwoo-hq/zwooc/pkg/model"
	"github.com/zwoo-hq/zwooc/pkg/tasks"
)

type customAdapter struct{}

var _ model.Adapter = (*customAdapter)(nil)

func NewAdapter() model.Adapter {
	return &customAdapter{}
}

func (a *customAdapter) CreateTask(c model.ProfileWrapper, extraArgs []string) tasks.Task {
	opts := c.GetOptions()
	data := helper.MapToStruct(opts, model.CustomOptions{})
	commandParts := strings.Split(data.Command, " ")

	cmd, additionalArgs := shared.CreateBaseCommand(commandParts[0], c, extraArgs)
	cmd.Args = append(cmd.Args, commandParts[1:]...)
	cmd.Args = append(cmd.Args, additionalArgs...)

	return tasks.NewCommandTask(c.GetName(), cmd)
}
