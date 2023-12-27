package main

import (
	"os"

	"github.com/urfave/cli/v2"
	"github.com/zwoo-hq/zwooc/pkg/config"
	"github.com/zwoo-hq/zwooc/pkg/helper"
	"github.com/zwoo-hq/zwooc/pkg/ui"
)

func createGlobalFlags() []cli.Flag {
	return []cli.Flag{
		&cli.BoolFlag{
			Name:    "no-tty",
			Aliases: []string{"s"},
			Usage:   "force disable tty features",
			Value:   false,
		},
		&cli.BoolFlag{
			Name:    "quite",
			Aliases: []string{"q"},
			Usage:   "disable all console output",
			Value:   false,
		},
		&cli.BoolFlag{
			Name:    "inline-output",
			Aliases: []string{"o"},
			Usage:   "inline output of tasks (in static mode)",
			Value:   false,
		},
		&cli.BoolFlag{
			Name:    "combine-output",
			Aliases: []string{"c"},
			Usage:   "combine output of tasks (in interactive mode)",
			Value:   false,
		},
		&cli.BoolFlag{
			Name:    "no-prefix",
			Aliases: []string{"p"},
			Usage:   "disable prefixing output of tasks with the task name",
			Value:   false,
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
	}
}

func execProfile(config config.Config, runMode string, c *cli.Context) error {
	viewOptions := ui.ViewOptions{
		DisableTUI:    c.Bool("no-tty"),
		QuiteMode:     c.Bool("quite"),
		InlineOutput:  c.Bool("inline-output"),
		CombineOutput: c.Bool("combine-output"),
		DisablePrefix: c.Bool("no-prefix"),
	}

	taskList, err := config.ResolveProfile(c.Args().First(), runMode)
	if err != nil {
		ui.HandleError(err)
	}

	ui.NewRunner(taskList, viewOptions)
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
	}
}

func execFragment(config config.Config, c *cli.Context) error {
	viewOptions := ui.ViewOptions{
		DisableTUI: c.Bool("no-tty"),
		QuiteMode:  c.Bool("quite"),
	}

	task, err := config.ResolvedFragment(c.Args().First())
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
		UseShortOptionHandling: true,
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
		},
	}

	if err := app.Run(os.Args); err != nil {
		ui.HandleError(err)
	}
}
