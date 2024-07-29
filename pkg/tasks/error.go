package tasks

import (
	"errors"
	"fmt"
	"strings"
)

var (
	ErrCancelled = errors.New("task cancelled")
)

type MultiTaskError struct {
	Errors map[string]error
}

func (mte MultiTaskError) Error() string {
	nodeIds := make([]string, 0, len(mte.Errors))
	for id := range mte.Errors {
		nodeIds = append(nodeIds, id)
	}
	return fmt.Sprintf("tasks %s failed", strings.Join(nodeIds, ", "))
}

func NewMultiTaskError(errors map[string]error) MultiTaskError {
	return MultiTaskError{
		Errors: errors,
	}
}
