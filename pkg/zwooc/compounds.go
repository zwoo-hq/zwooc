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
			completeFragments(conf)
		},
	}
}

func execCompound(conf config.Config, c *cli.Context) error {
	if c.Bool("dry-run") {
		return fmt.Errorf("--dry-run is currently not supported for compounds")
	}

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

	if isCI() && !c.Bool("no-ci") {
		viewOptions.DisableTUI = true
		viewOptions.InlineOutput = true
	}

	compoundKey := c.Args().First()
	compoundTasks, err := conf.LoadCompound(compoundKey)
	if err != nil {
		ui.HandleError(err)
	}

	ui.NewInteractiveRunner(compoundTasks, viewOptions, conf)
	return nil
}
