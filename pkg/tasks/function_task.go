package tasks

type functionTask struct {
	name    string
	execute func(cancel <-chan bool) error
}

func NewTask(name string, execute func(cancel <-chan bool) error) Task {
	return functionTask{
		name:    name,
		execute: execute,
	}
}

func (ft functionTask) Name() string {
	return ft.name
}

func (ft functionTask) Run(cancel <-chan bool) error {
	return ft.execute(cancel)
}
