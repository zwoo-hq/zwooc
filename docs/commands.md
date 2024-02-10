# zwooc Commands

This is an overview of the (wanted) zwooc functionality with implementation status.

- :white_check_mark: this feature is currently implemented
- :x: this feature is not currently implemented
- :question: implementation status unknown

## Projects

A project is a standalone application or library. The key of a project must equal the directory relative to the location of the config file. The project key may include `$COMP_WORDBREAKS` characters. Projects should define at least one profile.

| concept                  |       status       |
| ------------------------ | :----------------: |
| define projects          | :white_check_mark: |
| custom project directory |        :x:         |
| `vite-yarn` adapter      | :white_check_mark: |
| `dotnet` adapter         | :white_check_mark: |

## Profiles

A profile is a specific configuration in which a project can be run/built. The key of a profile shall not contain any `$COMP_WORDBREAKS` characters except colons `:` because these would break shell completion.

Profile shall contain at least one definition for one of the three available run modes `run`, `build` and `watch`.

The definitions may define further options in order to configure the environment in which the project is ran. These include `args` as of arguments which are passed to the process, `env` as of environment variables which will be available to the process, `alias` (soon to be `base`) as a reference to another profile of which the configuration will be inherited and `includeFragments` as of a list of fragments which will be run in parallel.

`args` are configured as an object with `key:value` pairs which will be translated into `--key value`. If the key already starts with a hyphen (`-`) the auto prefixing will be disabled. `env` values are passed as a list of strings in the format `VAR=value`. These value will be passed as is without any modification. Additionally adapters may include special env vars or arguments in order to achieve to output desired. Such special configuration will be used in order to enforce static or interactive mode or to provide special short hand configuration syntax.  

Furthermore definitions may include may include options which are depend on the adapter of the profile. These include `mode` for the `vite-yarn` adapter as a shorthand for the `--mode` arg. Profile definitions within `dotnet` adapter projects must contain an `project` option as of a reference to the desired `.csproj` file. 

| concept                               |       status       |
| ------------------------------------- | :----------------: |
| define profiles                       | :white_check_mark: |
| define `args` options                 | :white_check_mark: |
| don't enforce `--` prefix on args     |        :x:         |
| define env options                    | :white_check_mark: |
| define a base profile                 | :white_check_mark: |
| define included fragments             | :white_check_mark: |
| define `mode` in `vite-yarn` projects | :white_check_mark: |
| define `project` in `dotnet` projects | :white_check_mark: |

## Hooks
Any hook-able entity may define `$pre` and `$post` hooks. All profile definitions, fragments and compounds are considered hook-able. `$pre` hooks are always executed before the entity while `$post` hook are always executed after the entity. 

Hooks may define a command or reference a list of fragments. Due to fragments being hook-able themselves dependencies shall not be cyclic.

| concept                   |       status       |
| ------------------------- | :----------------: |
| define hooks              | :white_check_mark: |
| check if hooks are cyclic |        :x:         |

### Build Mode

| concept                     |       status       |
| --------------------------- | :----------------: |
| run a profile in build mode | :white_check_mark: |
| run build mode  interactive | :white_check_mark: |

### Run Mode

### Watch Mode

## Custom Tasks (Fragments)

Fragments are custom commands without any relation to profiles. Thus fragments can use and run any tool or commands they like. Fragments may be dependencies of profiles. Fragments may have dependencies in form of commands or other fragments on their own, these dependencies cant be cyclic.

The key of a fragment shall not contain any `$COMP_WORDBREAKS` characters except colons `:` because these would break shell completion.

Fragments must have at least one definition. The `$default` is executed whenever there is no more specific version to be found. Fragments may define more specific version of the command based on the current run mode or calling profile if its executed as a dependency. These can be defined via `<run mode>`,  `<profile>` or `<run mode>:<profile>`. When resolving a more specific version of a fragment the profile takes precedence over the run mode and run mode and profile takes precedence over one of them only.

When executing fragments via `exec` the it will execute the `$default` version, because no run mode or profile can be inferred. To allow executing a specific version of the fragment the run mode and or profile can be passed separated by a colon like `exec <fragment>:<run mode or profile>(:<profile>)`. 

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


## Utilities and options

Along the core functionality `zwooc` should provide additional utilities.

| concept                              |       status       |
| ------------------------------------ | :----------------: |
| version                              | :white_check_mark: |
| help                                 | :white_check_mark: |
| bash completion                      | :white_check_mark: |
| dependency/execution graph (dry run) |        :x:         |

Furthermore `zwooc` should provide global options in order to provide flexibility whilst executing tasks.

| concept                           |       status       |
| --------------------------------- | :----------------: |
| quite mode                        | :white_check_mark: |
| disable task output prefix        | :white_check_mark: |
| serial execution mode             | :white_check_mark: |
| set a max concurrency             | :white_check_mark: |
| loose (tolerant errors)           |        :x:         |
| skip hooks                        |        :x:         |
| exclude fragments                 |        :x:         |
| force disable tty                 | :white_check_mark: |
| inline output (static mode)       | :white_check_mark: |
| disable output (interactive mode) |        :x:         |
| combine output (interactive mode) | :white_check_mark: |
| no fullscreen (interactive mode)  |        :x:         |


# TODOs:

- :x: handle colons in bash completion better (https://stackoverflow.com/questions/10528695/how-to-reset-comp-wordbreaks-without-affecting-other-completion-script)
- :x: (BREAKING) rename alias to base
- :x: remove skip fragments from json config