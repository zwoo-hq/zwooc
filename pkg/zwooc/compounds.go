package zwooc

import (
	"fmt"

	"github.com/urfave/cli/v2"
	"github.com/zwoo-hq/zwooc/pkg/config"
	"github.com/zwoo-hq/zwooc/pkg/ui"
)

func CreateCompoundCommand() *cli.Command {
	return &cli.Command{
		Name:      "launch",
		Usage:     "launch a compound",
		ArgsUsage: "[compounds]",
		Flags:     CreateGlobalFlags(),
		Action: func(c *cli.Context) error {
			conf := loadConfig()
			return execCompound(conf, c)
		},
		BashComplete: func(c *cli.Context) {
			if c.NArg() > 0 {
				return
			}
			conf := loadConfig()
			completeCompounds(conf)
		},
	}
}

func execCompound(conf config.Config, c *cli.Context) error {
	if c.Bool("dry-run") {
		err := fmt.Errorf("--dry-run is currently not supported for compounds")
		ui.HandleError(err)
	}

	viewOptions := getViewOptions(c)
	ctx := config.NewContext(getLoadOptions(c, []string{}))
	compoundKey := c.Args().First()
	compoundTasks, err := conf.LoadCompound(compoundKey, ctx)
	if err != nil {
		ui.HandleError(err)
	}

	ui.NewInteractiveRunner(compoundTasks, viewOptions, conf)
	return nil
}
