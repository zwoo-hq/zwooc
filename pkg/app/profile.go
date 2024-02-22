package app

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

	viewOptions := ui.ViewOptions{
		DisableTUI:     c.Bool("no-tty"),
		QuiteMode:      c.Bool("quite"),
		InlineOutput:   c.Bool("inline-output"),
		CombineOutput:  c.Bool("combine-output"),
		DisablePrefix:  c.Bool("no-prefix"),
		MaxConcurrency: c.Int("max-concurrency"),
	}

	if c.Bool("serial") {
		viewOptions.MaxConcurrency = 1
	}

	if isCI() && !c.Bool("no-ci") {
		viewOptions.DisableTUI = true
		viewOptions.InlineOutput = true
	}

	args := c.Args().Tail()
	profileKey := c.Args().First()
	taskList, err := conf.ResolveProfile(profileKey, runMode, args)
	if err != nil {
		ui.HandleError(err)
	}

	list := *taskList.Flatten()
	list.RemoveEmptyStagesAndTasks()
	if runMode == config.ModeWatch || runMode == config.ModeRun {
		ui.NewInteractiveRunner(list, viewOptions, conf)
	} else {
		ui.NewRunner(list, viewOptions)
	}
	return nil
}
