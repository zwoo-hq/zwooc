package tasks

import "io"

func Empty() Task {
	return NewTask("noop", func(cancel <-chan bool, out io.Writer) error {
		return nil
	})
}
