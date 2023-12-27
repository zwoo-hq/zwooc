package tasks

import (
	"bytes"
	"io"
	"sync"
)

type CommandCapturer struct {
	data bytes.Buffer
	mu   sync.RWMutex
}

var _ io.Writer = (*CommandCapturer)(nil)

func NewCapturer() *CommandCapturer {
	return &CommandCapturer{}
}

func (cc *CommandCapturer) Write(p []byte) (n int, err error) {
	cc.mu.Lock()
	defer cc.mu.Unlock()
	return cc.data.Write(p)
}

func (cc *CommandCapturer) Bytes() []byte {
	cc.mu.RLock()
	defer cc.mu.RUnlock()
	return cc.data.Bytes()
}

func (cc *CommandCapturer) String() string {
	cc.mu.RLock()
	defer cc.mu.RUnlock()
	return cc.data.String()
}
