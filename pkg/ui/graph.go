package ui

import (
	"fmt"

	"github.com/zwoo-hq/zwooc/pkg/config"
	"github.com/zwoo-hq/zwooc/pkg/tasks"
)

func GraphDependencies(tasks *tasks.TaskTreeNode) {
	fmt.Printf("viewing %s\n", tasks.Name)
	tasks.RemoveEmptyNodes()
	printNode(tasks, "", true)
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
		newPrefix := "  │ "
		name := graphPreStyle.Render(config.KeyPre)
		info := graphInfoStyle.Render(fmt.Sprintf("(%d tasks)", len(node.Pre)))
		if len(node.Post) == 0 {
			newPrefix = "    "
			fmt.Printf("%s  └─┬%s %s\n", prefix, name, info)
		} else {
			fmt.Printf("%s  ├─┬%s %s\n", prefix, name, info)
		}

		for i, child := range node.Pre {
			printNode(child, prefix+newPrefix, i == len(node.Pre)-1)
		}
	}

	if len(node.Post) > 0 {
		fmt.Printf("%s  └─┬%s %s\n", prefix, graphPostStyle.Render(config.KeyPost), graphInfoStyle.Render(fmt.Sprintf("(%d tasks)", len(node.Post))))
		for i, child := range node.Post {
			printNode(child, prefix+"    ", i == len(node.Post)-1)
		}
	}
}
