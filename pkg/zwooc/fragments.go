package zwooc

import (
	"github.com/urfave/cli/v2"
	"github.com/zwoo-hq/zwooc/pkg/config"
	"github.com/zwoo-hq/zwooc/pkg/tasks"
	"github.com/zwoo-hq/zwooc/pkg/ui"
)

func CreateFragmentCommand() *cli.Command {
	return &cli.Command{
		Name:      "exec",
		Usage:     "execute a fragment",
		ArgsUsage: "[fragment] [extra arguments...]",
		Flags:     CreateGlobalFlags(),
		Action: func(c *cli.Context) error {
			conf := loadConfig()
			return execFragment(conf, c)
		},
		BashComplete: func(c *cli.Context) {
			if c.NArg() > 0 {
				return
			}
			conf := loadConfig()
			completeFragments(conf)
		},
	}
}

func execFragment(conf config.Config, c *cli.Context) error {
	if c.Bool("dry-run") {
		return graphTaskList(conf, c, "exec")
	}

	viewOptions := getViewOptions(c)
	ctx := config.NewContext(getLoadOptions(c, c.Args().Tail()))
	fragmentKey := c.Args().First()
	task, err := conf.LoadFragment(fragmentKey, ctx)
	if err != nil {
		ui.HandleError(err)
	}

	task.RemoveEmptyNodes()
	ui.NewRunner(tasks.NewCollection(task), viewOptions)
	return nil
}
