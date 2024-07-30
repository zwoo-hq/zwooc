package ui

import (
	"fmt"

	"github.com/zwoo-hq/zwooc/pkg/model"
	"github.com/zwoo-hq/zwooc/pkg/tasks"
)

type graphView struct {
	forest tasks.Collection
}

func GraphDependencies(collection tasks.Collection, name string) {
	fmt.Printf("%s - graphing dependency tree for %s\n", zwoocBranding, name)
	view := graphView{forest: collection}
	fmt.Println(view.View())
}

func (g *graphView) View() (s string) {
	for _, tasks := range g.forest {
		s += fmt.Sprintf("task %s ", graphHeaderStyle.Render(tasks.Name))
		s += graphInfoStyle.Render(fmt.Sprintf("(%d total linear equivalent stages)", tasks.CountStages())) + "\n"
		tasks.RemoveEmptyNodes()
		s += g.printGraphNode(tasks, "", true)
	}
	return
}

func (g *graphView) printGraphNode(node *tasks.TaskTreeNode, prefix string, isLast bool) (s string) {
	connector := "┬"
	if len(node.Pre) == 0 && len(node.Post) == 0 {
		connector = "─"
	}
	if isLast {
		s += fmt.Sprintf("%s└─%s%s %s\n", prefix, connector, graphMainStyle.Render(node.Name), graphInfoStyle.Render(node.Main.Name()))
	} else {
		s += fmt.Sprintf("%s├─%s%s %s\n", prefix, connector, graphMainStyle.Render(node.Name), graphInfoStyle.Render(node.Main.Name()))
	}

	if len(node.Pre) > 0 {
		newPrefix := "│ "
		prePrefix := "│ "
		if isLast {
			prePrefix = "  "
		}
		name := graphPreStyle.Render(model.KeyPre)
		info := graphInfoStyle.Render(fmt.Sprintf("(%d nodes)", len(node.Pre)))
		if len(node.Post) == 0 {
			newPrefix = "  "
			s += fmt.Sprintf("%s%s└─┬%s %s\n", prefix, prePrefix, name, info)
		} else {
			s += fmt.Sprintf("%s%s├─┬%s %s\n", prefix, prePrefix, name, info)
		}

		for i, child := range node.Pre {
			s += g.printGraphNode(child, prefix+prePrefix+newPrefix, i == len(node.Pre)-1)
		}
	}

	if len(node.Post) > 0 {
		postPrefix := "│ "
		if isLast {
			postPrefix = "  "
		}
		s += fmt.Sprintf("%s%s└─┬%s %s\n", prefix, postPrefix, graphPostStyle.Render(model.KeyPost), graphInfoStyle.Render(fmt.Sprintf("(%d tasks)", len(node.Post))))
		for i, child := range node.Post {
			s += g.printGraphNode(child, prefix+postPrefix+"  ", i == len(node.Post)-1)
		}
	}
	return
}
