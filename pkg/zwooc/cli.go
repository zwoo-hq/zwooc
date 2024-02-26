package zwooc

import (
	"fmt"
	"os"

	"github.com/urfave/cli/v2"
	"github.com/zwoo-hq/zwooc/pkg/config"
	"github.com/zwoo-hq/zwooc/pkg/helper"
	"github.com/zwoo-hq/zwooc/pkg/ui"
)

var (
	VERSION = "1.0.0-alpha.4"
)

var (
	CategoryStatic      = "Static mode (non TTY):"
	CategoryInteractive = "Interactive mode:"
	CategoryGeneral     = "General:"
	CategoryMisc        = "Miscellaneous:"
)

func loadConfig() config.Config {
	path, err := helper.FindFile("zwooc.config.json")
	if err != nil {
		ui.HandleError(err)
	}

	conf, err := config.Load(path)
	if err != nil {
		ui.HandleError(err)
	}
	return conf
}

func isCI() bool {
	return os.Getenv("CI") == "true"
}

func completeProfiles(c config.Config) {
	for _, profile := range c.GetProfiles() {
		if profile.Name() != config.KeyDefault {
			fmt.Print(profile.Name())
		}
	}
}

func completeFragments(c config.Config) {
	for _, fragment := range c.GetFragments() {
		if fragment.Name() != config.KeyDefault {
			fmt.Print(fragment.Name())
		}
	}
}

// func executeWithConfig(handler func(conf config.Config, c *cli.Context) error) func(c *cli.Context) error {
// 	return func(c *cli.Context) error {
// 		conf := loadConfig()
// 		return handler(conf, c)
// 	}
// }

func getLoadOptions(c *cli.Context, extraArgs []string) config.LoadOptions {
	return config.LoadOptions{
		SkipHooks: c.Bool("skip-hooks"),
		Exclude:   c.StringSlice("exclude"),
		ExtraArgs: extraArgs,
	}
}

func getViewOptions(c *cli.Context) ui.ViewOptions {
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
	return viewOptions
}
