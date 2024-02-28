package zwooc

import (
	"github.com/urfave/cli/v2"
	"github.com/zwoo-hq/zwooc/pkg/config"
	"github.com/zwoo-hq/zwooc/pkg/ui"
)

func CreateProfileCommand(mode, usage string) *cli.Command {
	return &cli.Command{
		Name:      mode,
		Usage:     usage,
		ArgsUsage: "[profile] [extra arguments...]",
		Flags:     CreateGlobalFlags(),
		Action: func(c *cli.Context) error {
			conf := loadConfig()
			return execProfile(conf, mode, c)
		},
		BashComplete: func(c *cli.Context) {
			if c.NArg() > 0 {
				return
			}
			conf := loadConfig()
			completeProfiles(conf)
		},
	}
}

func execProfile(conf config.Config, runMode string, c *cli.Context) error {
	if c.Bool("dry-run") {
		return graphTaskList(conf, c, runMode)
	}

	viewOptions := getViewOptions(c)
	ctx := config.NewContext(getLoadOptions(c, c.Args().Tail()))
	profileKey := c.Args().First()
	taskList, err := conf.LoadProfile(profileKey, runMode, ctx)
	if err != nil {
		ui.HandleError(err)
	}

	if runMode == config.ModeWatch || runMode == config.ModeRun || len(taskList) > 1 {
		ui.NewInteractiveRunner(taskList, viewOptions, conf)
	} else {
		list := taskList[0].Flatten()
		list.RemoveEmptyStagesAndTasks()
		ui.NewRunner(list, viewOptions)
	}
	return nil
}
