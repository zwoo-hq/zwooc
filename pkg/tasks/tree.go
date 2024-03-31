package tasks

import (
	"github.com/zwoo-hq/zwooc/pkg/helper"
)

// A TaskTreeNode represents a node in a task tree.
// A task tree is a tree structure that represents the order of tasks to be executed,
// each node has a main task and a collection of nodes that should be executed before
// and after the main task.
type TaskTreeNode struct {
	Name   string          // the name of the node
	Pre    []*TaskTreeNode // a collection of nodes that should be executed before the main task
	Main   Task            // the main task
	Post   []*TaskTreeNode // a collection of nodes that should be executed after the main task
	Parent *TaskTreeNode   // the parent node

	// IsLongRunning indicates whether the main task is long running
	// usually, only the main tasks of run or watch modes or such dependent fragments are long running
	IsLongRunning bool
}

func NewTaskTree(name string, mainTask Task, isLongRunning bool) *TaskTreeNode {
	return &TaskTreeNode{
		Name:          name,
		Pre:           []*TaskTreeNode{},
		Main:          mainTask,
		Post:          []*TaskTreeNode{},
		IsLongRunning: isLongRunning,
	}
}

// AddPreChild adds a child node to the pre collection.
// The child nodes parent is set to the current node.
func (t *TaskTreeNode) AddPreChild(child ...*TaskTreeNode) {
	for _, c := range child {
		c.Parent = t
	}
	t.Pre = append(t.Pre, child...)
}

// AddPostChild adds a child node to the post collection.
// The child nodes parent is set to the current node.
func (t *TaskTreeNode) AddPostChild(child ...*TaskTreeNode) {
	for _, c := range child {
		c.Parent = t
	}
	t.Post = append(t.Post, child...)
}

// FindNode returns a (child-)node with the given name.
func (t *TaskTreeNode) FindNode(name string) *TaskTreeNode {
	if t.Name == name {
		return t
	}
	for _, child := range helper.Concat(t.Pre, t.Post) {
		if node := child.FindNode(name); node != nil {
			return node
		}
	}
	return nil
}

// FindParent returns the parent node with the given name.
func (t *TaskTreeNode) FindParent(name string) *TaskTreeNode {
	parent := t
	for parent != nil {
		if parent.Name == name {
			return parent
		}
		parent = parent.Parent
	}

	return nil
}

// NodeID returns the unique identifier of the node in the current tree.
func (t *TaskTreeNode) NodeID() string {
	if t.Parent == nil {
		return t.Name
	}
	return helper.BuildName(t.Parent.NodeID(), t.Name)
}

// Iterate traverses the tree in depth-first order and calls the handler for each node.
func (t *TaskTreeNode) Iterate(handler func(node *TaskTreeNode)) {
	for _, pre := range t.Pre {
		pre.Iterate(handler)
	}
	handler(t)
	for _, post := range t.Post {
		post.Iterate(handler)
	}
}

// RemoveEmptyNodes removes all nodes that have an empty main task.
func (t *TaskTreeNode) RemoveEmptyNodes() {
	for i := 0; i < len(t.Pre); i++ {
		if IsEmptyTask(t.Pre[i].Main) {
			t.Pre = append(t.Pre[:i], t.Pre[i+1:]...)
			i--
		} else {
			t.Pre[i].RemoveEmptyNodes()
		}
	}
	for i := 0; i < len(t.Post); i++ {
		if IsEmptyTask(t.Post[i].Main) {
			t.Post = append(t.Post[:i], t.Post[i+1:]...)
			i--
		} else {
			t.Post[i].RemoveEmptyNodes()
		}
	}
}

// IsLinear returns true if the tree is linear, i.e. each node has at most one child node.
func (t *TaskTreeNode) IsLinear() bool {
	if len(t.Pre) > 1 || len(t.Post) > 1 {
		return false
	}

	for _, pre := range t.Pre {
		if !pre.IsLinear() {
			return false
		}
	}
	for _, post := range t.Post {
		if !post.IsLinear() {
			return false
		}
	}
	return true
}

// Flatten transforms the tree into a TaskList.
func (t *TaskTreeNode) Flatten() TaskList {
	list := NewTaskList(t.Name, []ExecutionStep{
		{
			Name:          t.Name,
			Tasks:         []Task{t.Main},
			IsLongRunning: t.IsLongRunning,
		},
	})

	preList := NewTaskList(helper.BuildName(t.Name, "pre"), []ExecutionStep{})
	for _, pre := range t.Pre {
		preList.MergePreAligned(pre.Flatten())
	}
	list.InsertBefore(preList)

	postList := NewTaskList(helper.BuildName(t.Name, "post"), []ExecutionStep{})
	for _, post := range t.Post {
		postList.MergePostAligned(post.Flatten())
	}
	list.InsertAfter(postList)

	return list
}

// CountStages returns the number of stages in the tree if it would be executed flattened.
func (t *TaskTreeNode) CountStages() int {
	preCount := 0
	for _, pre := range t.Pre {
		count := pre.CountStages()
		if count > preCount {
			preCount = count
		}
	}
	postCount := 0
	for _, post := range t.Post {
		count := post.CountStages()
		if count > postCount {
			postCount = count
		}
	}
	return 1 + preCount + postCount
}
