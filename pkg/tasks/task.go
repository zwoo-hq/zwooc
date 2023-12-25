package tasks

import "bytes"

type Task interface {
	Name() string
	Run(cancel <-chan bool) error
	Out() bytes.Buffer
}
