package zwooc

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
			return graphTaskList(conf, c, "")
		},
		BashComplete: func(c *cli.Context) {
			if c.NArg() > 1 {
				return
			}
			// complete first argument
			if c.NArg() == 0 {
				for _, mode := range []string{config.ModeBuild, config.ModeRun, config.ModeWatch, "exec"} {
					fmt.Println(mode)
				}
				return
			}

			conf := loadConfig()
			if c.Args().First() == "exec" {
				completeFragments(conf)
				return
			}
			completeProfiles(conf)
		},
	}
}

func graphTaskList(conf config.Config, c *cli.Context, defaultMode string) error {
	mode := c.Args().First()
	target := c.Args().Get(1)
	if defaultMode != "" {
		mode = defaultMode
		target = c.Args().First()
	}

	var tree *tasks.TaskTreeNode
	var err error

	if mode == "exec" {
		tree, err = conf.LoadFragment(target, []string{})
	} else if mode == "run" || mode == "watch" || mode == "build" {
		tree, err = conf.LoadProfile(target, mode, []string{})
	} else {
		err = fmt.Errorf("invalid mode: %s", mode)
	}

	if err != nil {
		ui.HandleError(err)
	}
	tree.RemoveEmptyNodes()
	ui.GraphDependencies(tree)
	return nil
}
