# zwooc concepts

This is an overview of the desired zwooc functionality with its implementation status.

- :white_check_mark: this feature is currently implemented
- :x: this feature is not currently implemented
- :question: implementation status unknown

## Projects

A project is a standalone application or library. The key of a project must equal the directory relative to the location of the config file. The project key may include `$COMP_WORDBREAKS` characters. Projects shall define at least one profile.

| concept                  |       status       |
| ------------------------ | :----------------: |
| define projects          | :white_check_mark: |
| custom project directory |        :x:         |
| `vite-yarn` adapter      | :white_check_mark: |
| `dotnet` adapter         | :white_check_mark: |

## Profiles

A profile is a specific set of parameters in which a project can be run/built. The key of a profile shall not contain any `$COMP_WORDBREAKS` characters except colons `:` because these would break shell completion.

Profile shall contain at least one definition for one of the three available run modes `run`, `build` and `watch`.

The definitions may define further options in order to configure the environment in which the project is run. These include `args` as of arguments which are passed to the process, `env` as of environment variables which will be available to the process, `alias` (soon to be `base`) as a reference to another profile of which the configuration will be inherited and `includeFragments` as of a list of fragments which will be run in parallel.

`args` are configured as an object with `key:value` pairs, which will be translated into `--key value`. If the key already starts with a hyphen (`-`) the auto prefixing will be disabled. `env` values are passed as a list of strings in the format `VAR=value`. These value will be passed as is without any modification. Additionally, adapters may include special env vars or arguments in order to achieve the output desired. Such special configuration will be used in order to enforce static or interactive mode or to provide special shorthand configuration syntax.  

Furthermore, definitions may include options which are dependent on the adapter of the profile. These include `mode` for the `vite-yarn` adapter as a shorthand for the `--mode` argument. Profile definitions within `dotnet` adapter projects must contain an `project` option as of a reference to the desired `.csproj` file. 

| concept                                |       status       |
| -------------------------------------- | :----------------: |
| define profiles                        | :white_check_mark: |
| define `args` options                  | :white_check_mark: |
| don't enforce `--` prefix on arguments |        :x:         |
| define env options                     | :white_check_mark: |
| define a base profile                  | :white_check_mark: |
| define included fragments              | :white_check_mark: |
| define `mode` in `vite-yarn` projects  | :white_check_mark: |
| define `project` in `dotnet` projects  | :white_check_mark: |

## Hooks
Any hook-able entity may define `$pre` and `$post` hooks. All profile definitions, fragments and compounds are considered hook-able. `$pre` hooks are always executed before the entity, while `$post` hook are always executed after the entity. 

Hooks may define a command or reference a list of fragments. Due to fragments being hook-able themselves, dependencies shall not be cyclic.

| concept                   |       status       |
| ------------------------- | :----------------: |
| define hooks              | :white_check_mark: |
| check if hooks are cyclic |        :x:         |

### Build Mode

The `build` run mode is defined for creating a compiled artifact from the code without executing the application or library. A standard profile definition is for the key `build` is required in order to execute a profile in build mode.

When executing a profile in build mode, a simpler task runner UI is used.

| concept                          |       status       |
| -------------------------------- | :----------------: |
| execute build mode               | :white_check_mark: |
| execute build mode (interactive) | :white_check_mark: |
| execute hooks                    | :white_check_mark: |
| execute included fragments       |     :question:     |

### Run & Watch Mode

The `run` run mode is defined for running an application once without applying changes after the source files were updated. 

The `watch` run mode is defined for running an application whilst watching the source files for changes and applying then seamlessly.

A standard profile definition is for the key `run` or `watch` is required in order to execute a profile in run or watch mode.

When executing a profile in run or watch mode, a more complex and feature-fuller runner is used when using an interactive runner.

| concept                                 |       status       |
| --------------------------------------- | :----------------: |
| execute run mode                        | :white_check_mark: |
| execute run mode (interactive)          | :white_check_mark: |
| execute watch mode                      | :white_check_mark: |
| execute watch mode (interactive)        | :white_check_mark: |
| execute hooks                           | :white_check_mark: |
| execute included fragments              |     :question:     |
| seamlessly switch between run and watch |        :x:         |

## Custom Tasks (Fragments)

Fragments are custom commands without any relation to profiles. Thus, fragments can use and run any tool or commands they like. Fragments may be dependencies of profiles. Fragments may have dependencies in the form of commands or other fragments on their own, these dependencies can't be cyclic.

The key of a fragment shall not contain any `$COMP_WORDBREAKS` characters except colons `:` because these would break shell completion.

Fragments must have at least one definition. The `$default` is executed whenever there is no more specific version to be found. Fragments may define a more specific version of the command based on the current run mode or calling profile if it's executed as a dependency. These can be defined via `<run mode>`,  `<profile>` or `<run mode>:<profile>`. When resolving a more specific version of a fragment, the profile takes precedence over the run mode and run mode and profile takes precedence over one of them only.

When executing fragments via `exec` will execute the `$default` version, because no run mode or profile can be inferred. To allow executing a specific version of the fragment, the run mode and or profile can be passed separated by a colon like `exec <fragment>:<run mode or profile>(:<profile>)`. 

| concept                                        |       status       |
| ---------------------------------------------- | :----------------: |
| define project scoped fragments                | :white_check_mark: |
| define global scoped fragments                 |     :question:     |
| execute fragments                              | :white_check_mark: |
| execute fragments (interactive)                | :white_check_mark: |
| pass extra arguments                           |     :question:     |
| execute with command dependencies              |        :x:         |
| execute with fragment dependencies             |        :x:         |
| detect cyclic dependencies                     |        :x:         |
| specific version based on run mode             | :white_check_mark: |
| specific version based on profile              |        :x:         |
| specific version based on run mode and profile | :white_check_mark: |

## Compounds

A compound is a collection of profiles or fragments which are executed in parallel. Compound tasks are executed in parallel by definition, this means other options like `max-concurrency` or `serial` do not take effect. 

The key of a compound shall not contain any `$COMP_WORDBREAKS` characters except colons `:` because these would break shell completion.

Compounds shall contain at least one profile or fragments. To configure a compound an object `profiles` with the profile key as the key and run desired run mode as value. Other options from profiles, like `base` and `includeFragments` apply here too.

| concept                         |       status       |
| ------------------------------- | :----------------: |
| define compounds                | :white_check_mark: |
| execute compounds               |        :x:         |
| execute compounds (interactive) |        :x:         |
| execute compounds (interactive) |        :x:         |
| execute hooks                   |        :x:         |
| execute included fragments      |        :x:         |

## Utilities and options

Along the core functionality, `zwooc` should provide additional utilities.

| concept                              |       status       |
| ------------------------------------ | :----------------: |
| version                              | :white_check_mark: |
| help                                 | :white_check_mark: |
| bash completion                      | :white_check_mark: |
| dependency/execution graph (dry run) |        :x:         |

Furthermore, `zwooc` should provide global options in order to provide flexibility whilst executing tasks.

| concept                           |       status       |
| --------------------------------- | :----------------: |
| quite mode                        | :white_check_mark: |
| disable task output prefix        | :white_check_mark: |
| serial execution mode             | :white_check_mark: |
| set a max concurrency             | :white_check_mark: |
| loose (tolerant errors)           |        :x:         |
| skip hooks                        |        :x:         |
| exclude fragments                 |        :x:         |
| force disable TTY                 | :white_check_mark: |
| inline output (static mode)       | :white_check_mark: |
| disable output (interactive mode) |        :x:         |
| combine output (interactive mode) | :white_check_mark: |
| no full screen (interactive mode) |        :x:         |

The hearth of the command line tool is the interactive task runner for TTYs. The interactive task runner shall allow for the execution of multiple tasks in parallel while scheduling all dependencies accordingly. The dynamic allocation of new tasks shall be able at runtime.

Since this UI is currently under heavy construction the following table only outlines future want-to-have features.

| concept                            |   status   |
| ---------------------------------- | :--------: |
| run tasks in parallel              | :question: |
| run decencies and hooks            | :question: |
| show task output in tabs           | :question: |
| combine the output of tasks        | :question: |
| allow scheduling tasks dynamically | :question: |
| kill tasks                         | :question: |
| handle errors in tasks             | :question: |


# TODOs:

- :x: handle colons in bash completion better (https://stackoverflow.com/questions/10528695/how-to-reset-comp-wordbreaks-without-affecting-other-completion-script)
- :x: (BREAKING) rename alias to base