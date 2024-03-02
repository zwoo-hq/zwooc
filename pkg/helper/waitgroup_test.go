package helper

import (
	"sync"
	"testing"
	"time"
)

func TestWaitFor(t *testing.T) {
	t.Run("should return channel that sends value once completed", func(tt *testing.T) {
		wg := &sync.WaitGroup{}
		target := &sync.WaitGroup{}
		wg.Add(2)
		target.Add(1)
		go func() {
			<-time.After(500 * time.Millisecond)
			target.Done()
			wg.Done()
		}()
		go func() {
			defer wg.Done()
			select {
			case <-WaitFor(target):
				tt.Log("WaitFor() returned")
				return
			case <-time.After(2 * time.Second):
				tt.Errorf("WaitFor() did not return after 1 second")
			}
		}()
		wg.Wait()
	})

}
