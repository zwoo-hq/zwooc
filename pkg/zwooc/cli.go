package zwooc

import (
	"fmt"
	"os"

	"github.com/zwoo-hq/zwooc/pkg/config"
	"github.com/zwoo-hq/zwooc/pkg/helper"
	"github.com/zwoo-hq/zwooc/pkg/ui"
)

var (
	VERSION = "1.0.0-alpha.3"
)

var (
	CategoryStatic      = "Static mode (non TTY):"
	CategoryInteractive = "Interactive mode:"
	CategoryGeneral     = "General:"
	CategoryFragments   = "Fragments:"
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
