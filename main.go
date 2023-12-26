package main

import (
	"os"

	"github.com/urfave/cli/v2"
	"github.com/zwoo-hq/zwooc/pkg/config"
	"github.com/zwoo-hq/zwooc/pkg/helper"
	"github.com/zwoo-hq/zwooc/pkg/ui"
)

type (
	Foo struct {
		Bar string `json:"bar"`
	}
)

func main() {
	path, err := helper.FindFile("zwoo.config.json")
	if err != nil {
		ui.HandleError(err)
	}

	conf, err := config.Load(path)
	if err != nil {
		ui.HandleError(err)
	}

	app := &cli.App{
		Name:  "zwooc",
		Usage: "the official cli for building and developing zwoo",
		Flags: []cli.Flag{
			&cli.BoolFlag{
				Name:  "no-tty",
				Usage: "force disable tty features",
				Value: false,
			},
			&cli.BoolFlag{
				Name:    "quite",
				Aliases: []string{"q"},
				Usage:   "disable all console output",
				Value:   false,
			},
		},
		Commands: []*cli.Command{
			{
				Name:  config.ModeRun,
				Usage: "run a profile",
				Action: func(c *cli.Context) error {
					return execProfile(conf, config.ModeRun, c)
				},
			},
			{
				Name:  config.ModeWatch,
				Usage: "run a profile with live reload enabled",
				Action: func(c *cli.Context) error {
					return execProfile(conf, config.ModeWatch, c)
				},
			},
			{
				Name:  config.ModeBuild,
				Usage: "build a profile",
				Action: func(c *cli.Context) error {
					return execProfile(conf, config.ModeBuild, c)
				},
			},
			{
				Name:  "exec",
				Usage: "execute a fragment",
				Action: func(c *cli.Context) error {
					return nil
				},
			},
			{
				Name:  "launch",
				Usage: "launch a compound",
				Action: func(c *cli.Context) error {
					return nil
				},
			},
		},
	}

	if err := app.Run(os.Args); err != nil {
		ui.HandleError(err)
	}
}

func execProfile(config config.Config, runMode string, c *cli.Context) error {
	viewOptions := ui.ViewOptions{
		DisableTUI: c.Bool("no-tty"),
		QuiteMode:  c.Bool("quite"),
	}

	taskList, err := config.ResolveProfile(c.Args().First(), runMode)
	if err != nil {
		ui.HandleError(err)
	}

	ui.NewRunner(taskList, viewOptions)
	return nil
}
