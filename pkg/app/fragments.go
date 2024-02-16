package app

import (
	"fmt"

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
			for _, fragment := range conf.GetFragments() {
				fmt.Println(fragment.Name())
			}
		},
	}
}

func execFragment(config config.Config, c *cli.Context) error {
	viewOptions := ui.ViewOptions{
		DisableTUI:     c.Bool("no-tty"),
		QuiteMode:      c.Bool("quite"),
		MaxConcurrency: c.Int("max-concurrency"),
	}

	if c.Bool("serial") {
		viewOptions.MaxConcurrency = 1
	}

	args := c.Args().Tail()
	fragmentKey := c.Args().First()
	task, err := config.ResolvedFragment(fragmentKey, args)
	if err != nil {
		ui.HandleError(err)
	}

	ui.NewFragmentRunner(task, viewOptions)
	return nil
}
