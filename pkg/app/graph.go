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
	var tree *tasks.TaskTreeNode
	var err error

	if mode == "exec" {
		tree, err = conf.ResolvedFragment(target, []string{})
	} else if mode == "run" || mode == "watch" || mode == "build" {
		tree, err = conf.ResolveProfile(target, mode, []string{})
	} else {
		err = fmt.Errorf("invalid mode: %s", mode)
	}

	if err != nil {
		ui.HandleError(err)
	}
	ui.GraphDependencies(tree)
	return nil
}
