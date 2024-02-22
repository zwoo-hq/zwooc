package app

import (
	"github.com/urfave/cli/v2"
	"github.com/zwoo-hq/zwooc/pkg/config"
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

func execFragment(config config.Config, c *cli.Context) error {
	if c.Bool("dry-run") {
		return graphTaskList(config, c, "exec")
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
	fragmentKey := c.Args().First()
	task, err := config.ResolvedFragment(fragmentKey, args)
	if err != nil {
		ui.HandleError(err)
	}

	list := task.Flatten()
	list.RemoveEmptyStagesAndTasks()
	ui.NewRunner(*list, viewOptions)
	return nil
}
