package main

import (
	"os"

	"github.com/urfave/cli/v2"
	"github.com/zwoo-hq/zwoo-builder/pkg/config"
	"github.com/zwoo-hq/zwoo-builder/pkg/helper"
	"github.com/zwoo-hq/zwoo-builder/pkg/ui"
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
	ui.Logger.Debugf("loading config file: %s\n", path)

	_, err = config.Load(path)
	if err != nil {
		ui.HandleError(err)
	}

	app := &cli.App{
		Name:  "%BINARYNAME%",
		Usage: "the official cli for building and developing zwoo",
		Commands: []*cli.Command{
			{
				Name:  config.ModeRun,
				Usage: "run a profile",
				Action: func(c *cli.Context) error {
					return nil
				},
			},
			{
				Name:  config.ModeWatch,
				Usage: "run a profile with live reload enabled",
				Action: func(c *cli.Context) error {
					return nil
				},
			},
			{
				Name:  config.ModeBuild,
				Usage: "build a profile",
				Action: func(c *cli.Context) error {
					return nil
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
