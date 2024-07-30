package zwooc

import (
	"github.com/urfave/cli/v2"
	"github.com/zwoo-hq/zwooc/pkg/config"
	"github.com/zwoo-hq/zwooc/pkg/tasks"
	"github.com/zwoo-hq/zwooc/pkg/ui"
	legacyui "github.com/zwoo-hq/zwooc/pkg/ui/legacy"
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
	if isDryRun(c) {
		return graphTaskTree(conf, c, "exec")
	}

	runnerOptions := getRunnerOptions(c)
	ctx := config.NewContext(getLoadOptions(c, c.Args().Tail()))
	fragmentKey := c.Args().First()
	task, err := conf.LoadFragment(fragmentKey, ctx)
	if err != nil {
		ui.HandleError(err)
	}
	task.RemoveEmptyNodes()

	if runnerOptions.UseLegacyRunner {
		viewOptions := getLegacyViewOptions(c)
		legacyui.NewRunner(task.Flatten(), viewOptions)
	} else {
		viewOptions := getViewOptions(c)
		provider := createRunner(tasks.NewCollection(task), runnerOptions)
		ui.NewRunner(tasks.NewCollection(task), provider, viewOptions)
	}
	return nil
}
