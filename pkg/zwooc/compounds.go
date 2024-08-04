package zwooc

import (
	"github.com/urfave/cli/v2"
	"github.com/zwoo-hq/zwooc/pkg/config"
	"github.com/zwoo-hq/zwooc/pkg/ui"
	legacyui "github.com/zwoo-hq/zwooc/pkg/ui/legacy"
)

func CreateCompoundCommand() *cli.Command {
	return &cli.Command{
		Name:      "launch",
		Usage:     "launch a compound",
		ArgsUsage: "[compounds]",
		Flags:     CreateGlobalFlags(),
		Action: func(c *cli.Context) error {
			conf := loadConfig()
			return execCompound(conf, c)
		},
		BashComplete: func(c *cli.Context) {
			if c.NArg() > 0 {
				return
			}
			conf := loadConfig()
			completeCompounds(conf)
		},
	}
}

func execCompound(conf config.Config, c *cli.Context) error {
	if c.Bool("dry-run") {
		return graphTaskTree(conf, c, "launch")
	}

	runnerOptions := getRunnerOptions(c)
	ctx := config.NewContext(getLoadOptions(c, []string{}))
	compoundKey := c.Args().First()
	compoundTasks, err := conf.LoadCompound(compoundKey, ctx)
	if err != nil {
		ui.HandleError(err)
	}

	if runnerOptions.UseLegacyRunner {
		viewOptions := getLegacyViewOptions(c)
		legacyui.NewInteractiveRunner(compoundTasks, viewOptions, conf)
	}

	viewOptions := getViewOptions(c)
	adapter := newStatusAdapter(compoundTasks, runnerOptions)
	ui.NewInteractiveView(compoundTasks, adapter.scheduler, viewOptions)
	return nil
}
