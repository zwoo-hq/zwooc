package tasks

type Task interface {
	Name() string
	Run(cancel <-chan bool) error
}
