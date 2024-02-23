package zwooc

import "github.com/urfave/cli/v2"

func CreateGlobalFlags() []cli.Flag {
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

		// Other
		&cli.BoolFlag{
			Name:     "no-ci",
			Usage:    "disable ci mode",
			Value:    false,
			Category: CategoryMisc,
		},
		&cli.BoolFlag{
			Name:     "dry-run",
			Usage:    "dry run mode, no tasks will be executed (same as graph)",
			Value:    false,
			Category: CategoryMisc,
		},
	}
}
