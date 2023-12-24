package helper

import (
	"os"
	"os/exec"
)

func NewCommand(cmd string, dir string) *exec.Cmd {
	if dir != "" {
		cmd := exec.Command("sh", "-c", cmd)
		cmd.Dir = dir
		cmd.Stdout = os.Stdout
		return cmd
	}
	return exec.Command("sh", "-c", cmd)
}
