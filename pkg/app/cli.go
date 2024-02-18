package app

import (
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
