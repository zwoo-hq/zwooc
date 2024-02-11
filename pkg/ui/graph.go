package ui

import (
	"fmt"

	"github.com/zwoo-hq/zwooc/pkg/tasks"
)

func GraphDependencies(tasks tasks.TaskList) {
	fmt.Printf("viewing %s\n", tasks.Name)
	for _, task := range tasks.Steps {
		fmt.Println("- " + task.Name)
		for _, part := range task.Tasks {
			fmt.Println("  - " + part.Name())
		}
	}
}
