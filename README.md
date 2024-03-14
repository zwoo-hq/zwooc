![zwooc-social-image](https://github.com/zwoo-hq/zwooc/assets/47701374/68f176e6-2eba-4abf-9e3a-a9d3a2d93619)

# zwooc

🚀 The official meta build tool for [zwoo](https://github.com/fabiankachlock/zwoo)!

`zwooc` is a **_meta_** built tool, which means it only leverages and orchestrates exiting build tools in order to produce outputs. zwooc aims to unify and simplify build tool configuration tailored to the use cases of zwoo.

---

## Installation

You can install zwooc from the 
- [GitHub Releases](https://github.com/zwoo-hq/zwooc/releases) or 
- via `go install github.com/zwoo-hq/zwooc/cmd/zwooc@latest` (go 1.21 is needed).

### Setting up auto completion

zwooc currently supports auto completion for `bash` & `zsh`, just add:

`source <(zwooc complete-bash)` into you `.bashrc` or

`source <(zwooc complete-zsh)` into your `.zshrc`

On bash completion on windows (for example in the git bash) will also provide completions for `zwooc.exe` 

## Usage

### Get Started

`$ zwooc init` will initialize a new zwooc workspace. This creates a new `zwooc.config.json` with some example content. 

### Basics

Your every day use commands will be:

`$ zwooc build|run|watch <key>` - execute a profile in the given run mode.

`$ zwooc exec <key>` - executes a fragment

`$ zwooc launch <key>` - launch a compound configuration.

Often used options are:

`$ zwooc exec -to` `-t` will disable TTY mode and `-o` will enable command output in static mode. These options are enable in CI by default.

If you want to pass some extra arguments to an command you can do this always behind the key, like:

`$ zwooc run dev --host` - in this case `dev` being a `vite-x` profile this will expose your dev server to the local network.

The drawback of this it, that all arguments targeted at zwooc must be passed before the key of the configuration to execute:

```diff
- zwooc run dev -to // does not work
+ zwooc run -to dev // does work
```

### Debugging Configuration

zwooc provides a handy tool for debugging the configuration.

`$ zwooc graph exec|build|run|watch|launch <key>` will print a tree will all tasks and their dependencies into the terminal. Adding the `--dry-run` flag to on of those commands will do the same.


### More Information

`$ zwooc h|help|-h|--help` prints an overview of all available commands including a short description.

`$ zwooc -v|--version` will print the version of zwooc.

## Using The Interactive Runner

> [!WARNING]  
> The interactive runner is still in early stage of development and may contain bugs.

The interactive runner as a TUI for running zwooc tasks.

It currently supports:
- a help view by pressing `h`
- a full screen view by pressing `f`
- multi tabs command output view
- switching tabs via `tab` `shift+tab` or mouse click (yes it has mouse support!)
- status indicator for pre and post tasks
- `esc` will close the full screen or help view
- `q` or `ctrl+c` will stop the runner gracefully (running al post tasks) pressing it a second time will cancel all running post tasks


## Concepts & Configuration

This is a rough overview about all concepts, for a full documentation see [`docs/concept.md`](https://github.com/zwoo-hq/zwooc/blob/main/docs/concept.md)

### Run Mode

There a 3 available run modes: `run`, `watch`, `build`

### Projects

Define a sub-project with an adapter. The adapter will handle how commands are build. A project contains a number of profiles which can be run. The name of the project must equal the subpath.

Available adapters are
- `vite-yarn`, `vite-npm`, `vite-pnpm`
- `tauri-yarn`, `tauri-npm`, `tauri-pnpm`
- `dotnet`
- `custom`

#### Profiles

A profile is a run configuration for running a project in a certain run mode.

### Fragments

Fragments are individual commands that can be run before/with/after profiles. They are not bound to the adapter und are run with the folder of the project they are defined in. Fragments can adapt teh current run mode and profile.

### Compounds

Compounds are a combination of profiles that can be run together. They are defined in root scope.

### Genral Concepts

All runnable entities can define pre and post actions via `$pre` and `$post`

## Example Configuration

```json
{
    "project1": {
        "$adapter": "vite",
        "$fragments": {
            "fragment1": {
                "$default": "" // always run
            },
            "fragment2": {
                "run": "",
                "watch": "",
                "build": ""
            },
            "fragment3": {
                "profile1": "", // run when its an dependency of profile1 in any mode
                "build:profile1": "" // run when its an dependency of profile1 in build mode
            },
            "fragment4": {
                "$pre": {}, // fragments can have pre & post hooks
                "$post": {}
            }
        },
        "profile1": {
            "run": true, // use the default command
            "watch": "", // use a custom command (TODO)
            "build": false // since build is set to false set, there is no build command
        },
        "profile2": {
            "build": {
                "args": {
                    "foo": "bar", // add --foo bar to command
                    "-foo": "bar", // add -foo bar to command
                    "--foo": "bar" // add -foo bar to command
                },
                "env": [
                    "FOO=BAR" // add FOO env var with bar as value
                ],
                "base": "profile1", // use profile 1 as base (and apply these configs)
                "skipFragments": true, // ignore all depended fragments
                "$pre": {
                    "fragments": [], // list of fragments to run before
                    "profiles": {}, // list of profiles with the mode to run before
                    "command": "" // command to run before
                },
                "$post": {
                    "fragments": [], // list of fragments to run after
                    "command": "" // command to run after
                }
            },
            "watch": {
                "includeFragments": [
                    "fragment2" // executes fragment2:watch:profile2 (or fragment2:watch)in parallel
                ]
            }
        },
    },
    "project2": {
        "$adapter": "dotnet",
        "profile3": {
            "project": "foo.bar.csproj" // pronet project must define a csproj
        }
    },
    "$fragments": {
        "foo": {} // fragment run in root folder
    },
    "$compounds": {
        "all": { // define a compound named all
            "profiles": {
                "profile1": "watch", // profile 1 should be started inw watch mode
                "profile2": "build" // profile 1 should be started inw watch mode
            },
            // ... all other base options ($pre, $post, skipFragments, includeFragments, base)
        }
    }
}
```

## Whishlist

- BuiltIn support for build-dir cleaning
- BuiltIn support for output copying
