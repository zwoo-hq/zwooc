package tasks

func Empty() Task {
	return NewTask("noop", func(cancel <-chan bool) error {
		return nil
	})
}
