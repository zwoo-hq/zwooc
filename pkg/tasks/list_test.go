package tasks

import (
	"testing"
)

func createTestStep(name string, tasks ...Task) ExecutionStep {
	return NewExecutionStep(name, tasks, false)
}

func createTestList(name string, steps ...ExecutionStep) TaskList {
	return NewTaskList(name, steps)
}

func TestMergePreAligned(t *testing.T) {
	listA := createTestList("a", createTestStep("a1", &emptyTask{"a1t"}), createTestStep("a2", &emptyTask{"a2t"}))
	listB := createTestList("b", createTestStep("b1", &emptyTask{"b1t"}), createTestStep("b2", &emptyTask{"b2t"}), createTestStep("b3", &emptyTask{"b3t"}))

	listA.MergePreAligned(listB)

	if len(listA.Steps) != 3 {
		t.Errorf("Expected 3 steps, got %d", len(listA.Steps))
	}
	if len(listA.Steps[0].Tasks) != 1 {
		t.Errorf("Expected step 0 to have 1 task, got %d", len(listA.Steps[0].Tasks))
	}
	if len(listA.Steps[1].Tasks) != 2 {
		t.Errorf("Expected step 1 to have 2 task, got %d", len(listA.Steps[1].Tasks))
	}
	if len(listA.Steps[2].Tasks) != 2 {
		t.Errorf("Expected step 2 to have 2 tasks, got %d", len(listA.Steps[2].Tasks))
	}
}

func TestMergePostAligned(t *testing.T) {
	listA := createTestList("a", createTestStep("a1", &emptyTask{"a1t"}), createTestStep("a2", &emptyTask{"a2t"}))
	listB := createTestList("b", createTestStep("b1", &emptyTask{"b1t"}), createTestStep("b2", &emptyTask{"b2t"}), createTestStep("b3", &emptyTask{"b3t"}))

	listA.MergePostAligned(listB)

	if len(listA.Steps) != 3 {
		t.Errorf("Expected 3 steps, got %d", len(listA.Steps))
	}
	if len(listA.Steps[0].Tasks) != 2 {
		t.Errorf("Expected step 0 to have 2 tasks, got %d", len(listA.Steps[0].Tasks))
	}
	if len(listA.Steps[1].Tasks) != 2 {
		t.Errorf("Expected step 1 to have 2 task, got %d", len(listA.Steps[1].Tasks))
	}
	if len(listA.Steps[2].Tasks) != 1 {
		t.Errorf("Expected step 2 to have 1 tasks, got %d", len(listA.Steps[2].Tasks))
	}
}
