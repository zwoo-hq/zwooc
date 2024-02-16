package app

import (
	"fmt"

	"github.com/urfave/cli/v2"
	"github.com/zwoo-hq/zwooc/pkg/config"
	"github.com/zwoo-hq/zwooc/pkg/ui"
)

func CreateProfileCommand(mode, usage string) *cli.Command {
	return &cli.Command{
		Name:      mode,
		Usage:     usage,
		ArgsUsage: "[profile]",
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
			for _, profile := range conf.GetProfiles() {
				fmt.Println(profile.Name())
			}
		},
	}
}

func execProfile(conf config.Config, runMode string, c *cli.Context) error {
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

	taskList, err := conf.ResolveProfile(c.Args().First(), runMode)
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
