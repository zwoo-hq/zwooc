package main

import (
	"embed"
	"fmt"
	"io"
	"os"
	"sort"

	"github.com/urfave/cli/v2"
	"github.com/zwoo-hq/zwooc/pkg/config"
	"github.com/zwoo-hq/zwooc/pkg/helper"
	"github.com/zwoo-hq/zwooc/pkg/ui"
)

var (
	CategoryStatic      = "Static mode (non TTY):"
	CategoryInteractive = "Interactive mode:"
	CategoryGeneral     = "General:"
	CategoryFragments   = "Fragments:"
)

//go:embed autocomplete/*
var autocompletion embed.FS

func createGlobalFlags() []cli.Flag {
	return []cli.Flag{
		// global
		&cli.BoolFlag{
			Name:     "quite",
			Aliases:  []string{"q"},
			Usage:    "disable all console output",
			Value:    false,
			Category: CategoryGeneral,
		},
		&cli.BoolFlag{
			Name:     "no-prefix",
			Aliases:  []string{"p"},
			Usage:    "disable prefixing output of tasks with the task name",
			Value:    false,
			Category: CategoryGeneral,
		},
		&cli.BoolFlag{
			Name:     "serial",
			Aliases:  []string{"s"},
			Usage:    "run tasks in serial instead of parallel",
			Value:    false,
			Category: CategoryGeneral,
		},
		&cli.IntFlag{
			Name:     "max-concurrency",
			Aliases:  []string{"c"},
			Usage:    "limit the max amount of parallel tasks",
			Category: CategoryGeneral,
		},
		&cli.BoolFlag{
			// TODO: implement
			Name:     "loose",
			Aliases:  []string{"l"},
			Usage:    "ignores errors in tasks and continues",
			Value:    false,
			Category: CategoryGeneral,
		},
		&cli.BoolFlag{
			// TODO: implement
			Name:     "skip-hooks",
			Aliases:  []string{"n"},
			Usage:    "ignore all $pre and $post hooks",
			Value:    false,
			Category: CategoryGeneral,
		},

		// Static mode
		&cli.BoolFlag{
			Name:     "no-tty",
			Aliases:  []string{"t"},
			Usage:    "force disable tty features",
			Value:    false,
			Category: CategoryStatic,
		},
		&cli.BoolFlag{
			Name:     "inline-output",
			Aliases:  []string{"o"},
			Usage:    "inline output of tasks in static mode",
			Value:    false,
			Category: CategoryStatic,
		},

		// Interactive mode
		&cli.BoolFlag{
			// TODO: implement
			Name: "no-output",
			// Aliases:  []string{"o"},
			Usage:    "disable command output capturing in interactive mode",
			Value:    false,
			Category: CategoryInteractive,
		},
		&cli.BoolFlag{
			// TODO: implement
			Name: "combine-output",
			// Aliases:  []string{"c"},
			Usage:    "combine output of tasks in interactive mode",
			Value:    false,
			Category: CategoryInteractive,
		},
		&cli.BoolFlag{
			// TODO: implement
			Name:     "no-fullscreen",
			Aliases:  []string{"i"},
			Usage:    "inlines the interactive view ",
			Value:    false,
			Category: CategoryInteractive,
		},

		// Fragments
		&cli.StringSliceFlag{
			// TODO: implement
			Name:     "exclude",
			Aliases:  []string{"e"},
			Usage:    "excludes certain fragments from being executed",
			Category: CategoryFragments,
		},
	}
}

func createProfileCommand(mode, usage string, conf config.Config) *cli.Command {
	return &cli.Command{
		Name:      mode,
		Usage:     usage,
		ArgsUsage: "[profile]",
		Flags:     createGlobalFlags(),
		Action: func(c *cli.Context) error {
			return execProfile(conf, mode, c)
		},
		BashComplete: func(c *cli.Context) {
			if c.NArg() > 0 {
				return
			}
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

	if runMode == config.ModeWatch || runMode == config.ModeRun {
		ui.NewInteractiveRunner(taskList, viewOptions, conf)
	} else {
		ui.NewRunner(taskList, viewOptions)
	}
	return nil
}

func createFragmentCommand(conf config.Config) *cli.Command {
	return &cli.Command{
		Name:      "exec",
		Usage:     "execute a fragment",
		ArgsUsage: "[fragment] [extra arguments...]",
		Flags:     createGlobalFlags(),
		Action: func(c *cli.Context) error {
			return execFragment(conf, c)
		},
		BashComplete: func(c *cli.Context) {
			if c.NArg() > 0 {
				return
			}
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
		Name:                   "zwooc",
		Usage:                  "the official cli for building and developing zwoo",
		Flags:                  createGlobalFlags(),
		Suggest:                true,
		UseShortOptionHandling: true,
		EnableBashCompletion:   true,
		Commands: []*cli.Command{
			createProfileCommand(config.ModeRun, "run a profile", conf),
			createProfileCommand(config.ModeWatch, "run a profile with live reload enabled", conf),
			createProfileCommand(config.ModeBuild, "build a profile", conf),
			createFragmentCommand(conf),
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
