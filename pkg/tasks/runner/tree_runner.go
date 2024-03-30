package runner

import (
	"fmt"
	"runtime"
	"sync"

	"github.com/zwoo-hq/zwooc/pkg/tasks"
)

type TreeStatusNode struct {
	Name      string
	Status    TaskStatus
	PreNodes  []*TreeStatusNode
	PostNodes []*TreeStatusNode
	Parent    *TreeStatusNode
}

type TaskTreeRunner struct {
	root           *tasks.TaskTreeNode
	runningNodes   []*tasks.TaskTreeNode
	status         *TreeStatusNode
	updates        chan int
	cancel         chan bool
	cancelComplete chan error
	mutex          sync.RWMutex
	maxConcurrency int
}

func NewTaskTreeRunner(root *tasks.TaskTreeNode, maxConcurrency int) *TaskTreeRunner {
	ticketAmount := maxConcurrency
	if ticketAmount < 1 {
		ticketAmount = runtime.NumCPU()
	}
	status := buildStatus(root)

	return &TaskTreeRunner{
		root:           root,
		runningNodes:   []*tasks.TaskTreeNode{},
		status:         status,
		updates:        make(chan int, 1),
		cancel:         make(chan bool),
		cancelComplete: make(chan error),
		maxConcurrency: ticketAmount,
		mutex:          sync.RWMutex{},
	}
}

func buildStatus(root *tasks.TaskTreeNode) *TreeStatusNode {
	status := &TreeStatusNode{
		Name:      root.Name,
		Status:    StatusPending,
		PreNodes:  []*TreeStatusNode{},
		PostNodes: []*TreeStatusNode{},
	}

	for _, pre := range root.Pre {
		status.PreNodes = append(status.PreNodes, buildStatus(pre))
	}

	for _, post := range root.Post {
		status.PostNodes = append(status.PostNodes, buildStatus(post))
	}

	return status
}

func findStatus(status *TreeStatusNode, target *tasks.TaskTreeNode) *TreeStatusNode {
	path := []string{target.Name}
	for target.Parent != nil {
		path = append([]string{target.Parent.Name}, path...)
		target = target.Parent
	}
	fmt.Println(path)
	current := status
	for _, name := range path {
		if current.Name == name {
			return current
		}

		for _, pre := range current.PreNodes {
			if pre.Name == name {
				current = pre
				continue
			}
		}

		for _, post := range current.PostNodes {
			if post.Name == name {
				current = post
				continue
			}
		}
		break
	}
	return nil
}
