package config

import (
	"fmt"
	"strings"

	"github.com/zwoo-hq/zwooc/pkg/helper"
	"github.com/zwoo-hq/zwooc/pkg/tasks"
)

func (c Config) ResolvedFragment(key string, extraArgs []string) (tasks.Task, error) {
	parts := strings.Split(key, ":")
	mode := ""
	profile := ""

	if len(parts) >= 2 {
		key = parts[0]
		mode = parts[1]
	}
	if len(parts) >= 3 {
		key = parts[2]
	}

	fragment, err := c.resolveFragment(key, mode, profile)
	if err != nil {
		return tasks.Empty(), err
	}

	return fragment.GetTask(extraArgs), nil
}

func (c Config) resolveFragment(key, mode, profile string) (ResolvedFragment, error) {
	target, found := helper.FindBy(c.fragments, func(f Fragment) bool {
		return f.Name() == key
	})
	if !found {
		return ResolvedFragment{}, fmt.Errorf("fragment '%s' not found", key)
	}

	return target.GetConfig(mode, profile)
}
