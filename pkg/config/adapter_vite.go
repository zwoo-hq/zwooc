package config

import (
	"os"
	"os/exec"
	"path/filepath"

	"github.com/zwoo-hq/zwooc/pkg/tasks"
)

func CreateViteTask(c ResolvedProfile) tasks.Task {
	cmd := exec.Command("yarn")
	profileOptions := c.GetProfileOptions()

	cmd.Env = append(cmd.Env, profileOptions.Env...)

	for k, v := range profileOptions.Args {
		cmd.Args = append(cmd.Args, "--"+k)
		cmd.Args = append(cmd.Args, v)
	}

	if wd, err := os.Getwd(); err == nil {
		cmd.Dir = filepath.Join(wd, c.Directory)
	}

	return tasks.NewCommandTask(c.Name, cmd)
}
