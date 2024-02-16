package app

import (
	"fmt"

	"github.com/urfave/cli/v2"
	"github.com/zwoo-hq/zwooc/pkg/config"
	"github.com/zwoo-hq/zwooc/pkg/tasks"
	"github.com/zwoo-hq/zwooc/pkg/ui"
)

func CreateGraphCommand() *cli.Command {
	return &cli.Command{
		Name:      "graph",
		Usage:     "display a graph of tasks",
		ArgsUsage: "[run|watch|build|exec] [profile or fragment]",
		Flags:     CreateGlobalFlags(),
		Action: func(c *cli.Context) error {
			conf := loadConfig()
			return graphTaskList(conf, c)
		},
		BashComplete: func(c *cli.Context) {
			// TODO: implement
		},
	}
}

func graphTaskList(conf config.Config, c *cli.Context) error {
	mode := c.Args().First()
	target := c.Args().Get(1)
	var taskList tasks.TaskList
	var err error

	if mode == "exec" {
		var task tasks.Task
		task, err = conf.ResolvedFragment(target, []string{})
		taskList = tasks.TaskList{
			Name: "fragment",
			Steps: []tasks.ExecutionStep{
				{
					Name:  "fragment",
					Tasks: []tasks.Task{task},
				},
			},
		}
	} else if mode == "run" || mode == "watch" || mode == "build" {
		var tree *tasks.TaskTreeNode
		tree, err = conf.ResolveProfile(target, mode)
		if tree != nil {
			taskList = *tree.Flatten()
		}
	} else {
		err = fmt.Errorf("invalid mode: %s", mode)
	}

	if err != nil {
		ui.HandleError(err)
	}
	taskList.RemoveEmptyStagesAndTasks()
	ui.GraphDependencies(taskList)
	return nil
}
