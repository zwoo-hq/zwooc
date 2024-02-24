package config

import (
	"fmt"
	"strings"

	"github.com/zwoo-hq/zwooc/pkg/helper"
	"github.com/zwoo-hq/zwooc/pkg/tasks"
)

func normalizeFragmentKey(fullKey string) (key, mode, profile string) {
	parts := strings.Split(fullKey, ":")
	key = parts[0]
	if len(parts) >= 2 {
		mode = parts[1]
	}
	if len(parts) >= 3 {
		profile = parts[2]
	}
	return
}

func (c Config) LoadFragment(key string, extraArgs []string) (*tasks.TaskTreeNode, error) {
	key, mode, profile := normalizeFragmentKey(key)
	fragment, err := c.resolveFragment(key, mode, profile)
	if err != nil {
		return nil, err
	}

	node := tasks.NewTaskTree(fragment.Name, fragment.GetTask(extraArgs), false)
	err = c.loadAllHooks(fragment, node, mode, profile)
	if err != nil {
		return nil, err
	}
	return node, nil
}

func (c Config) resolveFragment(key, mode, profile string) (ResolvedFragment, error) {
	target, found := helper.FindBy(c.fragments, func(f Fragment) bool {
		return f.Name() == key
	})
	if !found {
		return ResolvedFragment{}, fmt.Errorf("fragment '%s' not found", key)
	}

	return target.ResolveConfig(mode, profile)
}
