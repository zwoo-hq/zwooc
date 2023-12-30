package ui

import (
	"bytes"
	"io"
)

type notifyWriter struct {
	buf     *bytes.Buffer
	updates chan string
}

var _ io.Writer = (*notifyWriter)(nil)

func NewNotifyWriter() *notifyWriter {
	return &notifyWriter{
		buf:     bytes.NewBuffer(nil),
		updates: make(chan string, 8),
	}
}

func (w *notifyWriter) Write(p []byte) (n int, err error) {
	n, err = w.buf.Write(p)
	w.updates <- w.buf.String()
	return n, err
}

func (w *notifyWriter) String() string {
	return w.buf.String()
}

func (w *notifyWriter) Close() error {
	close(w.updates)
	return nil
}
