package main

import (
	"embed"
	"fmt"
	"io"
	"os"
	"sort"

	"github.com/urfave/cli/v2"
	"github.com/zwoo-hq/zwooc/pkg/config"
	"github.com/zwoo-hq/zwooc/pkg/ui"
	"github.com/zwoo-hq/zwooc/pkg/zwooc"
)

//go:embed autocomplete/*
var autocompletion embed.FS

func main() {
	cli.VersionPrinter = func(c *cli.Context) {
		fmt.Println(c.App.Version)
	}

	zwooc := &cli.App{
		Name:    "zwooc",
		Usage:   "the official cli for building and developing zwoo",
		Version: zwooc.VERSION,

		Flags:                  zwooc.CreateGlobalFlags(),
		Suggest:                true,
		UseShortOptionHandling: true,
		EnableBashCompletion:   true,
		Commands: []*cli.Command{
			zwooc.CreateProfileCommand(config.ModeRun, "run a profile"),
			zwooc.CreateProfileCommand(config.ModeWatch, "run a profile with live reload enabled"),
			zwooc.CreateProfileCommand(config.ModeBuild, "build a profile"),
			zwooc.CreateFragmentCommand(),
			zwooc.CreateCompoundCommand(),
			zwooc.CreateGraphCommand(),
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

	sort.Sort(cli.FlagsByName(zwooc.Flags))
	sort.Sort(cli.CommandsByName(zwooc.Commands))

	if err := zwooc.Run(os.Args); err != nil {
		ui.HandleError(err)
	}
}
