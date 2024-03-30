package runner

type TaskStatus int

const (
	StatusPending TaskStatus = iota
	StatusRunning
	StatusDone
	StatusError
	StatusCanceled
)
