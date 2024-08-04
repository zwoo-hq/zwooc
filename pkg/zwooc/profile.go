package zwooc

import (
	"github.com/urfave/cli/v2"
	"github.com/zwoo-hq/zwooc/pkg/config"
	"github.com/zwoo-hq/zwooc/pkg/model"
	"github.com/zwoo-hq/zwooc/pkg/ui"
	legacyui "github.com/zwoo-hq/zwooc/pkg/ui/legacy"
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
		return graphTaskTree(conf, c, runMode)
	}

	runnerOptions := getRunnerOptions(c)
	ctx := config.NewContext(getLoadOptions(c, c.Args().Tail()))
	profileKey := c.Args().First()
	allTasks, err := conf.LoadProfile(profileKey, runMode, ctx)
	if err != nil {
		ui.HandleError(err)
	}

	for _, task := range allTasks {
		task.RemoveEmptyNodes()
	}

	if runnerOptions.UseLegacyRunner {
		viewOptions := getLegacyViewOptions(c)
		if runMode == model.ModeWatch || runMode == model.ModeRun || len(allTasks) > 1 {
			legacyui.NewInteractiveRunner(allTasks, viewOptions, conf)
		} else {
			legacyui.NewRunner(allTasks[0].Flatten(), viewOptions)
		}
		return nil
	}

	viewOptions := getViewOptions(c)
	if runMode == model.ModeWatch || runMode == model.ModeRun || len(allTasks) > 1 {
		adapter := newStatusAdapter(allTasks, runnerOptions)
		ui.NewInteractiveView(allTasks, adapter.scheduler, viewOptions)
	} else {
		adapter := newStatusAdapter(allTasks, runnerOptions)
		ui.NewView(allTasks, adapter.scheduler.SimpleStatusProvider, viewOptions)
	}
	return nil
}
