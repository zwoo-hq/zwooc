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
	if len(parts) == 2 {
		mode = parts[1]
	} else if len(parts) == 3 {
		mode = parts[1]
		profile = parts[2]
	} else if len(parts) > 3 {
		key = strings.Join(parts[:len(parts)-2], ":")
		mode = parts[len(parts)-2]
		profile = parts[len(parts)-1]
	}
	return
}

func combineFragmentKey(key, mode, profile string) string {
	return fmt.Sprintf("%s:%s:%s", key, mode, profile)
}

func (c Config) LoadFragment(rawKey string, ctx loadingContext) (*tasks.TaskTreeNode, error) {
	key, mode, profile := normalizeFragmentKey(rawKey)
	if ctx.excludes(key) || ctx.excludes(rawKey) {
		return nil, ErrTargetExcluded
	}

	fragment, err := c.resolveFragment(key, mode, profile)
	if err != nil {
		// try with raw key
		fragment, err = c.resolveFragment(rawKey, "", "")
		if err != nil {
			return nil, err
		}
	}

	node := tasks.NewTaskTree(fragment.Name, fragment.GetTask(ctx.getArgs()), false)
	if !ctx.skipHooks {
		err = c.loadAllHooks(fragment, node, mode, profile, ctx.withCaller(fragment.Name))
		if err != nil {
			return nil, err
		}
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
