package helper

import "sync"

func WaitFor(wg *sync.WaitGroup) chan bool {
	done := make(chan bool)
	go func() {
		wg.Wait()
		close(done)
	}()
	return done
}
