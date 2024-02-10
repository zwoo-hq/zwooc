# zwooc Commands

This is an overview of the (wanted) zwooc functionality with implementation status.

- :white_check_mark: this feature is currently implemented
- :x: this feature is not currently implemented
- :question: implementation status unknown

# Build Mode

# Run Mode

# Watch Mode

# Custom Tasks (Fragments)

Fragments are custom commands without any relation to profiles. Thus fragments can use and run any tool or commands they like. Fragments may be dependencies of profiles. Fragments may have dependencies in form of commands or other fragments on their own, these dependencies cant be cyclic.

Fragments must have at least one definition. The `$default` is executed whenever there is no more specific version to be found. Fragments may define more specific version of the command based on the current run mode or calling profile if its executed as a dependency. These can be defined via `<run mode>`,  `<profile>` or `<run mode>:<profile>`. When resolving a more specific version of a fragment the profile takes precedence over the run mode and run mode and profile takes precedence over one of them only.

When executing fragments via `exec` the it will execute the `$default` version, because no run mode or profile can be inferred. To allow executing a specific version of the fragment the run mode and or profile can be passed separated by a colon like `exec <fragment>:<run mode or profile>(:<profile>)`. 

| concept                                        | status             |
| ---------------------------------------------- | ------------------ |
| execute fragments                              | :white_check_mark: |
| execute fragments (interactive)                | :white_check_mark: |
| execute with command dependencies              | :x:                |
| execute with fragment dependencies             | :x:                |
| detect cyclic dependencies                     | :x:                |
| specific version based on run mode             | :white_check_mark: |
| specific version based on profile              | :x:                |
| specific version based on run mode and profile | :white_check_mark: |


# Compounds

