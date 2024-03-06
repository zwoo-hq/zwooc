package config

import (
	"errors"
	"fmt"
	"strings"
)

var (
	ErrTargetExcluded     error = errors.New("target entity itself is excluded")
	ErrCircularDependency error = CircularDependencyError{}
)

type CircularDependencyError struct {
	target string
	caller []string
}

func (e CircularDependencyError) Error() string {
	return fmt.Sprintf("circular dependency detected: '%s' from %s", e.target, strings.Join(e.caller, " -> "))
}

func (e CircularDependencyError) Is(target error) bool {
	_, ok := target.(CircularDependencyError)
	return ok
}
