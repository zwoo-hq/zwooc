package tasks

import (
	"bytes"
	"io"
)

type CommandPrefixer struct {
	prefix    []byte
	isMidLine bool
	dest      io.Writer
}

var _ io.Writer = (*CommandPrefixer)(nil)

func NewPrefixer(prefix string, dest io.Writer) *CommandPrefixer {
	return &CommandPrefixer{
		prefix: []byte(prefix),
		dest:   dest,
	}
}

func (r *CommandPrefixer) Write(p []byte) (n int, err error) {
	// split bytes into lines
	lines := bytes.Split(p, []byte("\n"))
	for i, line := range lines {
		if len(line) == 0 {
			continue
		}

		if i == 0 && r.isMidLine {
			// do not prepend prefix
			r.isMidLine = false
		} else {
			line = append(r.prefix, line...)
			line = append(line, []byte("\n")...)
		}

		if i == len(lines)-1 {
			// if the last line does not end with \n
			// we need to remember that we are in the middle of a line
			r.isMidLine = !bytes.HasSuffix(line, []byte("\n"))
		}

		// write the line
		n, err = r.dest.Write(line)
		if err != nil {
			return n, err
		}
	}
	return len(p), err
}
