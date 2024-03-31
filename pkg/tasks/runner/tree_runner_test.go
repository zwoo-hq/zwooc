package runner

import (
	"testing"

	"github.com/zwoo-hq/zwooc/pkg/helper"
	"github.com/zwoo-hq/zwooc/pkg/tasks"
)

func TestBuildStatus(t *testing.T) {
	t.Run("builds status tree from tasks", func(t *testing.T) {
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

		if status.Name() != "root" {
			t.Errorf("Expected root, got %s", status.Name())
		}
		if status.Status() != StatusPending {
			t.Errorf("Expected pending, got %d", status.Status)
		}
		if len(status.PreNodes) != 1 {
			t.Errorf("Expected 1, got %d", len(status.PreNodes))
		}
		if len(status.PostNodes) != 1 {
			t.Errorf("Expected 1, got %d", len(status.PostNodes))
		}
		if status.PreNodes[0].Name() != "pre1" {
			t.Errorf("Expected pre1, got %s", status.PreNodes[0].Name())
		}
		if status.PostNodes[0].Name() != "post1" {
			t.Errorf("Expected post1, got %s", status.PostNodes[0].Name())
		}
		if len(status.PreNodes[0].PreNodes) != 1 {
			t.Errorf("Expected 1, got %d", len(status.PreNodes[0].PreNodes))
		}
		if len(status.PostNodes[0].PostNodes) != 1 {
			t.Errorf("Expected 1, got %d", len(status.PostNodes[0].PostNodes))
		}
		if status.PreNodes[0].PreNodes[0].Name() != "pre1-1" {
			t.Errorf("Expected pre1-1, got %s", status.PreNodes[0].PreNodes[0].Name())
		}
		if status.PostNodes[0].PostNodes[0].Name() != "post1-1" {
			t.Errorf("Expected post1-1, got %s", status.PostNodes[0].PostNodes[0].Name())
		}
	})
}

func TestFindStatus(t *testing.T) {
	t.Run("finds status node for task", func(t *testing.T) {
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
		if node.Name() != "root" {
			t.Errorf("Expected root, got %s", node.Name())
		}

		node = findStatus(status, tree.Pre[0])
		if node.Name() != "pre1" {
			t.Errorf("Expected pre1, got %s", node.Name())
		}

		node = findStatus(status, tree.Post[0])
		if node.Name() != "post1" {
			t.Errorf("Expected post1, got %s", node.Name())
		}

		node = findStatus(status, tree.Pre[0].Pre[0])
		if node.Name() != "pre1-1" {
			t.Errorf("Expected pre1-1, got %s", node.Name())
		}

		node = findStatus(status, tree.Post[0].Post[0])
		if node.Name() != "post1-1" {
			t.Errorf("Expected post1-1, got %s", node.Name())
		}
	})
}

func TestGetStartingNodes(t *testing.T) {
	t.Run("gets starting nodes from tree", func(t *testing.T) {
		tree := &tasks.TaskTreeNode{
			Name: "root",
			Pre: []*tasks.TaskTreeNode{
				{
					Name: "pre1",
				},
				{
					Name: "pre2",
					Pre: []*tasks.TaskTreeNode{
						{
							Name: "pre2-1",
						},
						{
							Name: "pre2-2",
						},
					},
				},
				{
					Name: "pre3",
					Pre: []*tasks.TaskTreeNode{
						{
							Name: "pre3-1",
						},
						{
							Name: "pre3-2",
							Pre: []*tasks.TaskTreeNode{
								{
									Name: "pre3-2-1",
								},
							},
						},
					},
				},
			},
		}

		startingNodes := getStartingNodes(tree)
		expected := []string{"pre1", "pre2-1", "pre2-2", "pre3-1", "pre3-2-1"}

		if len(expected) != len(startingNodes) {
			t.Errorf("Expected len %d, got %d", len(expected), len(startingNodes))
		}

		for _, expectedName := range expected {
			if _, found := helper.FindBy(startingNodes, func(node *tasks.TaskTreeNode) bool {
				return node.Name == expectedName
			}); !found {
				t.Errorf("Expected %s, got <not found>", expectedName)
			}
		}

	})
}
