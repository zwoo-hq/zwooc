package tasks

import (
	"bytes"
	"io"
)

type NotifyWriter struct {
	buf     *bytes.Buffer
	Updates chan string
}

var _ io.Writer = (*NotifyWriter)(nil)

func NewNotifyWriter() *NotifyWriter {
	return &NotifyWriter{
		buf:     bytes.NewBuffer(nil),
		Updates: make(chan string, 8),
	}
}

func (w *NotifyWriter) Write(p []byte) (n int, err error) {
	n, err = w.buf.Write(p)
	w.Updates <- w.buf.String()
	return n, err
}

func (w *NotifyWriter) String() string {
	return w.buf.String()
}

func (w *NotifyWriter) Close() error {
	close(w.Updates)
	return nil
}
