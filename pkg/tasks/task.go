package tasks

import (
	"io"
	"strings"

	"github.com/zwoo-hq/zwooc/pkg/helper"
)

type Task interface {
	Name() string
	Run(cancel <-chan bool) error
	Pipe(destination io.Writer)
}

type Collection []*TaskTreeNode

func NewCollection(nodes ...*TaskTreeNode) Collection {
	return Collection(nodes)
}

func (c Collection) GetName() string {
	if len(c) == 0 {
		return ""
	} else if len(c) == 1 {
		return c[0].Name
	} else {
		mapToName := func(node *TaskTreeNode) string {
			return node.Name
		}
		return strings.Join(helper.MapTo(c, mapToName), ", ")
	}
}
