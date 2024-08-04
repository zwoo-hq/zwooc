package runner

import "runtime"

type ConcurrencyProvider interface {
	Acquire() int
	Release(ticket int)
	Schedule(action func())
	Close()
}

type SharedConcurrencyProvider struct {
	maxConcurrency int
	tickets        chan int
}

var _ ConcurrencyProvider = (*SharedConcurrencyProvider)(nil)

func NewSharedProvider(maxConcurrency int) ConcurrencyProvider {
	if maxConcurrency < 1 {
		maxConcurrency = runtime.NumCPU()
	}

	provider := &SharedConcurrencyProvider{
		maxConcurrency: maxConcurrency,
		tickets:        make(chan int, maxConcurrency),
	}

	for i := 0; i < maxConcurrency; i++ {
		provider.tickets <- i
	}

	return provider
}

func (c *SharedConcurrencyProvider) Acquire() int {
	return <-c.tickets
}

func (c *SharedConcurrencyProvider) Release(ticket int) {
	c.tickets <- ticket
}

func (c *SharedConcurrencyProvider) Schedule(action func()) {
	ticket := c.Acquire()
	action()
	c.Release(ticket)
}

func (c *SharedConcurrencyProvider) Close() {
	close(c.tickets)
}
