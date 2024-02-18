package tasks

import (
	"fmt"
	"testing"
)

func TestTree(t *testing.T) {
	tree := NewTaskTree("root", Empty(), true)
	childChildA := NewTaskTree("childChildA", Empty(), false)
	childChildB := NewTaskTree("childChildB", Empty(), false)

	childA := NewTaskTree("childA", Empty(), false)
	childA.AddPreChild(childChildA)
	childA.AddPostChild(childChildB)

	childB := NewTaskTree("childB", Empty(), false)
	tree.AddPreChild(childA)
	tree.AddPostChild(childA)
	tree.AddPreChild(childB)
	tree.AddPostChild(childB)

	list := tree.Flatten()
	fmt.Println(list)

	if len(list.Steps) != 7 {
		t.Errorf("Expected 7 steps, got %d", len(list.Steps))
	}
}
