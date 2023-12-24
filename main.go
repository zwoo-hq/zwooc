package main

import (
	"os"
	"time"

	"github.com/urfave/cli/v2"
	"github.com/zwoo-hq/zwooc/pkg/config"
	"github.com/zwoo-hq/zwooc/pkg/helper"
	"github.com/zwoo-hq/zwooc/pkg/tasks"
	"github.com/zwoo-hq/zwooc/pkg/ui"
)

type (
	Foo struct {
		Bar string `json:"bar"`
	}
)

func sleepTask(name string, dur time.Duration) tasks.Task {
	return tasks.NewTask(name, func(cancel <-chan bool) error {
		select {
		case <-cancel:
			return nil
		case <-time.After(dur * time.Second):
			return nil
		}
	})
}

func main() {
	// ui.CreateFullScreenView(config.TaskList{
	// 	Name: "test",
	// 	Steps: []config.ExecutionStep{
	// 		{
	// 			Name:        "pre",
	// 			RunParallel: false,
	// 			Tasks: []tasks.Task{
	// 				sleepTask("task 1", 2),
	// 				sleepTask("task 2", 3),
	// 				sleepTask("task 3", 4),
	// 			},
	// 		},
	// 		{
	// 			Name: "main",
	// 			Tasks: []tasks.Task{
	// 				sleepTask("dotnet build", 7),
	// 			},
	// 		},
	// 		{
	// 			Name: "post",
	// 			Tasks: []tasks.Task{
	// 				sleepTask("cleanup", 4),
	// 			},
	// 		},
	// 	},
	// })

	// ui.PrintSuccess("test", 2*time.Second)

	path, err := helper.FindFile("zwoo.config.json")
	if err != nil {
		ui.HandleError(err)
	}
	ui.Logger.Debugf("loading config file: %s", path)

	conf, err := config.Load(path)
	if err != nil {
		ui.HandleError(err)
	}

	app := &cli.App{
		Name:  "zwooc",
		Usage: "the official cli for building and developing zwoo",
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

	taskList, err := config.ResolveProfile(c.Args().First(), runMode)
	if err != nil {
		ui.HandleError(err)
	}

	ui.RunStatic(taskList)
	return nil
}
