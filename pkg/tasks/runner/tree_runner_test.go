package runner

import (
	"testing"

	"github.com/zwoo-hq/zwooc/pkg/tasks"
)

func TestBuildStatus(t *testing.T) {
	t.Run("TestbuildStatus", func(t *testing.T) {
		tree := &tasks.TaskTreeNode{
			Name: "root",
			Pre: []*tasks.TaskTreeNode{
				{
					Name: "pre1",
					Pre: []*tasks.TaskTreeNode{
						{
							Name: "pre1-1",
						},
					},
				},
			},
			Post: []*tasks.TaskTreeNode{
				{
					Name: "post1",
					Post: []*tasks.TaskTreeNode{
						{
							Name: "post1-1",
						},
					},
				},
			},
		}

		status := buildStatus(tree)

		if status.Name != "root" {
			t.Errorf("Expected root, got %s", status.Name)
		}
		if status.Status != StatusPending {
			t.Errorf("Expected pending, got %d", status.Status)
		}
		if len(status.PreNodes) != 1 {
			t.Errorf("Expected 1, got %d", len(status.PreNodes))
		}
		if len(status.PostNodes) != 1 {
			t.Errorf("Expected 1, got %d", len(status.PostNodes))
		}
		if status.PreNodes[0].Name != "pre1" {
			t.Errorf("Expected pre1, got %s", status.PreNodes[0].Name)
		}
		if status.PostNodes[0].Name != "post1" {
			t.Errorf("Expected post1, got %s", status.PostNodes[0].Name)
		}
		if len(status.PreNodes[0].PreNodes) != 1 {
			t.Errorf("Expected 1, got %d", len(status.PreNodes[0].PreNodes))
		}
		if len(status.PostNodes[0].PostNodes) != 1 {
			t.Errorf("Expected 1, got %d", len(status.PostNodes[0].PostNodes))
		}
		if status.PreNodes[0].PreNodes[0].Name != "pre1-1" {
			t.Errorf("Expected pre1-1, got %s", status.PreNodes[0].PreNodes[0].Name)
		}
		if status.PostNodes[0].PostNodes[0].Name != "post1-1" {
			t.Errorf("Expected post1-1, got %s", status.PostNodes[0].PostNodes[0].Name)
		}
	})
}

func TestFindStatus(t *testing.T) {
	t.Run("TestFindStatus", func(t *testing.T) {
		tree := &tasks.TaskTreeNode{
			Name: "root",
			Pre: []*tasks.TaskTreeNode{
				{
					Name: "pre1",
					Pre: []*tasks.TaskTreeNode{
						{
							Name: "pre1-1",
						},
					},
				},
			},
			Post: []*tasks.TaskTreeNode{
				{
					Name: "post1",
					Post: []*tasks.TaskTreeNode{
						{
							Name: "post1-1",
						},
					},
				},
			},
		}
		tree.Pre[0].Parent = tree
		tree.Pre[0].Pre[0].Parent = tree.Pre[0]
		tree.Post[0].Parent = tree
		tree.Post[0].Post[0].Parent = tree.Post[0]

		status := buildStatus(tree)

		node := findStatus(status, tree)
		if node.Name != "root" {
			t.Errorf("Expected root, got %s", node.Name)
		}

		node = findStatus(status, tree.Pre[0])
		if node.Name != "pre1" {
			t.Errorf("Expected pre1, got %s", node.Name)
		}

		node = findStatus(status, tree.Post[0])
		if node.Name != "post1" {
			t.Errorf("Expected post1, got %s", node.Name)
		}

		node = findStatus(status, tree.Pre[0].Pre[0])
		if node.Name != "pre1-1" {
			t.Errorf("Expected pre1-1, got %s", node.Name)
		}

		node = findStatus(status, tree.Post[0].Post[0])
		if node.Name != "post1-1" {
			t.Errorf("Expected post1-1, got %s", node.Name)
		}
	})
}
