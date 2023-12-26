package tasks

import (
	"bytes"
	"io"
)

type CommandCapturer struct {
	data bytes.Buffer
}

var _ io.Writer = &CommandCapturer{}

func NewCapturer() *CommandCapturer {
	return &CommandCapturer{}
}

func (cc *CommandCapturer) Write(p []byte) (n int, err error) {
	return cc.data.Write(p)
}

func (cc *CommandCapturer) Bytes() []byte {
	return cc.data.Bytes()
}

func (cc *CommandCapturer) String() string {
	return cc.data.String()
}
