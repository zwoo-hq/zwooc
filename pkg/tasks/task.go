package tasks

import (
	"io"
)

type Task interface {
	Name() string
	Run(cancel <-chan bool) error
	Pipe(destination io.Writer)
}
