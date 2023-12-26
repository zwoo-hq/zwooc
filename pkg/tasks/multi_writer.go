package tasks

import "io"

type multiWriter struct {
	writers []io.Writer
}

var _ io.Writer = (*multiWriter)(nil)

func newMultiWriter(writers ...io.Writer) *multiWriter {
	allWriters := []io.Writer{}
	for _, w := range writers {
		if mw, ok := w.(*multiWriter); ok {
			allWriters = append(allWriters, mw.writers...)
		} else {
			allWriters = append(allWriters, w)
		}
	}
	return &multiWriter{allWriters}
}

func (t *multiWriter) Write(p []byte) (n int, err error) {
	for _, w := range t.writers {
		n, err = w.Write(p)
		if err != nil {
			return
		}
		if n != len(p) {
			err = io.ErrShortWrite
			return
		}
	}
	return len(p), nil
}

func (t *multiWriter) Pipe(w io.Writer) {
	t.writers = append(t.writers, w)
}
