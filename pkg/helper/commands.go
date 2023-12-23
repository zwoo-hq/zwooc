package helper

import "os/exec"

func NewCommand(cmd string, dir string) *exec.Cmd {
	if dir != "" {
		cmd := exec.Command("sh", "-c", cmd)
		cmd.Dir = dir
		return cmd
	}
	return exec.Command("sh", "-c", cmd)
}
