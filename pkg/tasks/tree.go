package tasks

import (
	"github.com/zwoo-hq/zwooc/pkg/helper"
)

type TaskTreeNode struct {
	Name          string
	Pre           []*TaskTreeNode
	Main          Task
	Post          []*TaskTreeNode
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

func (t *TaskTreeNode) AddPreChild(child ...*TaskTreeNode) {
	t.Pre = append(t.Pre, child...)
}

func (t *TaskTreeNode) AddPostChild(child ...*TaskTreeNode) {
	t.Post = append(t.Post, child...)
}

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

func (t *TaskTreeNode) Flatten() *TaskList {
	list := NewTaskList(t.Name, []ExecutionStep{
		{
			Name:          t.Name,
			Tasks:         []Task{t.Main},
			IsLongRunning: t.IsLongRunning,
		},
	})

	preList := NewTaskList(helper.BuildName(t.Name, "pre"), []ExecutionStep{})
	for _, pre := range t.Pre {
		preList.MergePreAligned(*pre.Flatten())
	}
	list.InsertBefore(preList)

	postList := NewTaskList(helper.BuildName(t.Name, "post"), []ExecutionStep{})
	for _, post := range t.Post {
		postList.MergePostAligned(*post.Flatten())
	}
	list.InsertAfter(postList)

	return &list
}

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
