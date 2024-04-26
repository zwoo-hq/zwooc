package ui

import (
	"fmt"

	"github.com/zwoo-hq/zwooc/pkg/model"
	"github.com/zwoo-hq/zwooc/pkg/tasks"
)

func GraphDependencies(collection tasks.Collection, name string) {
	fmt.Printf("%s - graphing dependency tree for %s\n", zwoocBranding, name)
	for _, tasks := range collection {
		fmt.Printf("task %s ", graphHeaderStyle.Render(tasks.Name))
		fmt.Println(graphInfoStyle.Render(fmt.Sprintf("(%d total linear equivalent stages)", tasks.CountStages())))
		tasks.RemoveEmptyNodes()
		printNode(tasks, "", true)
	}
}

func printNode(node *tasks.TaskTreeNode, prefix string, isLast bool) {
	connector := "┬"
	if len(node.Pre) == 0 && len(node.Post) == 0 {
		connector = "─"
	}
	if isLast {
		fmt.Printf("%s└─%s%s %s\n", prefix, connector, graphMainStyle.Render(node.Name), graphInfoStyle.Render(node.Main.Name()))
	} else {
		fmt.Printf("%s├─%s%s %s\n", prefix, connector, graphMainStyle.Render(node.Name), graphInfoStyle.Render(node.Main.Name()))
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
			fmt.Printf("%s%s└─┬%s %s\n", prefix, prePrefix, name, info)
		} else {
			fmt.Printf("%s%s├─┬%s %s\n", prefix, prePrefix, name, info)
		}

		for i, child := range node.Pre {
			printNode(child, prefix+prePrefix+newPrefix, i == len(node.Pre)-1)
		}
	}

	if len(node.Post) > 0 {
		postPrefix := "│ "
		if isLast {
			postPrefix = "  "
		}
		fmt.Printf("%s%s└─┬%s %s\n", prefix, postPrefix, graphPostStyle.Render(model.KeyPost), graphInfoStyle.Render(fmt.Sprintf("(%d tasks)", len(node.Post))))
		for i, child := range node.Post {
			printNode(child, prefix+postPrefix+"  ", i == len(node.Post)-1)
		}
	}
}
