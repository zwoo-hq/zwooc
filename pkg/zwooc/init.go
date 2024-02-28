package zwooc

import (
	"embed"
	"fmt"
	"os"

	"github.com/urfave/cli/v2"
	"github.com/zwoo-hq/zwooc/pkg/ui"
)

//go:embed template.zwooc.config.json
var fs embed.FS

func CreateInitCommand() *cli.Command {
	return &cli.Command{
		Name:  "init",
		Usage: "init a new zwooc workspace",
		Action: func(c *cli.Context) error {
			data, err := fs.ReadFile("template.zwooc.config.json")
			if err != nil {
				ui.HandleError(err)
				return err
			}

			if _, err = os.Stat("zwooc.config.json"); err == nil {
				err = fmt.Errorf("zwooc.config.json already exists in the current directory")
				ui.HandleError(err)
				return err
			}

			err = os.WriteFile("zwooc.config.json", data, 0644)
			if err != nil {
				ui.HandleError(err)
				return err
			}
			ui.PrintSuccess("successfully created a zwooc.config.json in the current directory")
			return nil
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
