package main

import (
	"embed"
	"fmt"
	"io"
	"os"
	"sort"

	"github.com/urfave/cli/v2"
	"github.com/zwoo-hq/zwooc/pkg/app"
	"github.com/zwoo-hq/zwooc/pkg/config"
	"github.com/zwoo-hq/zwooc/pkg/ui"
)

//go:embed autocomplete/*
var autocompletion embed.FS

func main() {
	cli.VersionPrinter = func(c *cli.Context) {
		fmt.Println(c.App.Version)
	}

	app := &cli.App{
		Name:    "zwooc",
		Usage:   "the official cli for building and developing zwoo",
		Version: app.VERSION,

		Flags:                  app.CreateGlobalFlags(),
		Suggest:                true,
		UseShortOptionHandling: true,
		EnableBashCompletion:   true,
		Commands: []*cli.Command{
			app.CreateProfileCommand(config.ModeRun, "run a profile"),
			app.CreateProfileCommand(config.ModeWatch, "run a profile with live reload enabled"),
			app.CreateProfileCommand(config.ModeBuild, "build a profile"),
			app.CreateFragmentCommand(),
			app.CreateGraphCommand(),
			{
				Name:  "launch",
				Usage: "launch a compound",
				Action: func(c *cli.Context) error {
					return nil
				},
			},
			{
				// TODO: when cliv3 comes out this is no longer needed
				Name:  "completion-script",
				Usage: "generate shell completion script",
				Action: func(c *cli.Context) error {
					f, err := autocompletion.Open("autocomplete/bash_autocomplete")
					if err != nil {
						return err
					}

					content, err := io.ReadAll(f)
					if err != nil {
						return err
					}
					fmt.Println(string(content))
					return nil
				},
			},
		},
	}

	sort.Sort(cli.FlagsByName(app.Flags))
	sort.Sort(cli.CommandsByName(app.Commands))

	if err := app.Run(os.Args); err != nil {
		ui.HandleError(err)
	}
}
