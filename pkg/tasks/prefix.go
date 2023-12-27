package tasks

import (
	"io"
)

type CommandPrefixer struct {
	prefix []byte
	dest   io.Writer
}

var _ io.Writer = (*CommandPrefixer)(nil)

func NewPrefixer(prefix string, dest io.Writer) *CommandPrefixer {
	return &CommandPrefixer{
		prefix: []byte(prefix),
		dest:   dest,
	}
}

func (r *CommandPrefixer) Write(p []byte) (n int, err error) {
	n, err = r.dest.Write(append(r.prefix, p...))
	return n - len(r.prefix), err
}
