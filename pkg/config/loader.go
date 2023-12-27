package config

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"github.com/zwoo-hq/zwooc/pkg/helper"
	"github.com/zwoo-hq/zwooc/pkg/tasks"
)

type Config struct {
	baseDir string
	raw     map[string]interface{}
}

func Load(path string) (Config, error) {
	content, err := os.ReadFile(path)
	if err != nil {
		return Config{}, err
	}

	data := map[string]interface{}{}
	err = json.Unmarshal(content, &data)
	if err != nil {
		return Config{}, err
	}

	return Config{
		baseDir: filepath.Dir(path),
		raw:     data,
	}, nil
}

func (c Config) GetProfiles() ([]Profile, error) {
	profiles := []Profile{}

	for projectKey, projectValue := range c.raw {
		if !IsReservedKey(projectKey) {
			project := projectValue.(map[string]interface{})
			var projectAdapter string
			if adapter, ok := project[KeyAdapter]; ok {
				projectAdapter = adapter.(string)
			} else {
				return []Profile{}, fmt.Errorf("project '%s' is missing adapter", projectKey)
			}

			for profileKey, profileValue := range project {
				if !IsReservedKey(profileKey) {
					newProfile := Profile{
						name:      profileKey,
						adapter:   projectAdapter,
						directory: filepath.Join(c.baseDir, projectKey),
						raw:       profileValue.(map[string]interface{}),
					}
					profiles = append(profiles, newProfile)
				}
			}
		}
	}

	return profiles, nil
}

func (c Config) GetFragments() ([]Fragment, error) {
	fragments := []Fragment{}

	for projectKey, projectValue := range c.raw {
		if !IsReservedKey(projectKey) {
			project := projectValue.(map[string]interface{})
			if fragmentDefinitions, ok := project[KeyFragment]; ok {
				for fragmentKey, fragmentValue := range fragmentDefinitions.(map[string]interface{}) {
					newFragment := Fragment{
						name:      fragmentKey,
						directory: filepath.Join(c.baseDir, projectKey),
						raw:       fragmentValue.(map[string]interface{}),
					}
					fragments = append(fragments, newFragment)
				}
			}
		}
	}

	if fragmentDefinitions, ok := c.raw[KeyFragment]; ok {
		for fragmentKey, fragmentValue := range fragmentDefinitions.(map[string]interface{}) {
			newFragment := Fragment{
				name:      fragmentKey,
				directory: c.baseDir,
				raw:       fragmentValue.(map[string]interface{}),
			}
			fragments = append(fragments, newFragment)
		}
	}

	return fragments, nil
}

func (c Config) GetCompounds() ([]Compound, error) {
	compounds := []Compound{}

	if compoundDefinitions, ok := c.raw[KeyCompound]; ok {
		for compoundKey, compoundValue := range compoundDefinitions.(map[string]interface{}) {
			newCompound := newCompound(compoundKey, compoundValue.(map[string]interface{}))
			compounds = append(compounds, newCompound)
		}
	}
	return compounds, nil
}

func (c Config) ResolveProfile(key, mode string) (TaskList, error) {
	config, err := c.resolveRunConfig(key, mode)
	if err != nil {
		return TaskList{}, err
	}

	name := helper.BuildName(key, mode)
	preStage, err := c.resolveHook(KeyPre, config, config.GetPreHooks())
	if err != nil {
		return TaskList{}, err
	}

	postStage, err := c.resolveHook(KeyPost, config, config.GetPostHooks())
	if err != nil {
		return TaskList{}, err
	}

	mainTask, err := config.GetTask()
	if err != nil {
		return TaskList{}, err
	}

	steps := []ExecutionStep{}
	if len(preStage) > 0 {
		steps = append(steps, ExecutionStep{
			Tasks:       preStage,
			Name:        helper.BuildName(name, KeyPre),
			RunParallel: true,
		})
	}
	steps = append(steps, ExecutionStep{
		Tasks: []tasks.Task{mainTask},
		Name:  name,
	})
	if len(postStage) > 0 {
		steps = append(steps, ExecutionStep{
			Tasks:       postStage,
			Name:        helper.BuildName(name, KeyPost),
			RunParallel: true,
		})
	}

	return TaskList{
		Name:  name,
		Steps: steps,
	}, nil
}

func (c Config) resolveRunConfig(key, mode string) (ResolvedProfile, error) {
	if key == "" {
		key = KeyDefault
	}

	profiles, err := c.GetProfiles()
	if err != nil {
		return ResolvedProfile{}, err
	}
	target, found := helper.FindBy(profiles, func(p Profile) bool {
		return p.Name() == key
	})
	if !found {
		return ResolvedProfile{}, fmt.Errorf("profile '%s' not found", key)
	}

	config, err := target.GetConfig(mode)
	if err != nil {
		return ResolvedProfile{}, err
	}
	return config, nil
}

func (c Config) resolveHook(hookType string, profile ResolvedProfile, hook HookOptions) ([]tasks.Task, error) {
	baseName := hookType
	taskList := []tasks.Task{}
	if hook.Command != "" {
		taskList = append(taskList, tasks.NewBasicCommandTask(baseName, hook.Command, profile.Directory, []string{}))
	}

	for _, fragment := range hook.Fragments {
		fragmentConfig, err := c.resolveFragment(fragment, profile.Mode, profile.Name)
		if err != nil {
			return []tasks.Task{}, err
		}
		taskList = append(taskList, tasks.NewBasicCommandTask(helper.BuildName(baseName, fragment), fragmentConfig.Command, fragmentConfig.Directory, []string{}))
	}
	return taskList, nil
}
