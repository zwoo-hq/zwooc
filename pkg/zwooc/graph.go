package zwooc

import (
	"fmt"

	"github.com/urfave/cli/v2"
	"github.com/zwoo-hq/zwooc/pkg/config"
	"github.com/zwoo-hq/zwooc/pkg/model"
	"github.com/zwoo-hq/zwooc/pkg/tasks"
	"github.com/zwoo-hq/zwooc/pkg/ui"
)

func CreateGraphCommand() *cli.Command {
	return &cli.Command{
		Name:      "graph",
		Usage:     "display a graph of tasks",
		ArgsUsage: "[run|watch|build|exec|launch] [profile or fragment]",
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
				for _, mode := range []string{model.ModeBuild, model.ModeRun, model.ModeWatch, "exec"} {
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

	ctx := config.NewContext(getLoadOptions(c, []string{}))
	var forest tasks.Collection
	var err error

	if mode == "exec" {
		var task *tasks.TaskTreeNode
		task, err = conf.LoadFragment(target, ctx)
		forest = tasks.NewCollection(task)
	} else if mode == "launch" {
		forest, err = conf.LoadCompound(target, ctx)
	} else if mode == "run" || mode == "watch" || mode == "build" {
		forest, err = conf.LoadProfile(target, mode, ctx)
	} else {
		err = fmt.Errorf("invalid mode: %s", mode)
	}

	if err != nil {
		ui.HandleError(err)
	}
	for _, tree := range forest {
		tree.RemoveEmptyNodes()
	}
	ui.GraphDependencies(forest)
	return nil
}
