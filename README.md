![zwooc-social-image](https://github.com/zwoo-hq/zwooc/assets/47701374/68f176e6-2eba-4abf-9e3a-a9d3a2d93619)

# zwooc

The official build system for [zwoo](https://github.com/fabiankachlock/zwoo)!

---

## Install

You can install zwooc from the  [GitHub Releases](https://github.com/zwoo-hq/zwooc/releases) or via `go install github.com/zwoo-hq/zwooc/cmd/zwooc@latest` (go 1.20 is needed).

## Concepts

This is a rough overview about all concepts, for a full documentation see [`docs/concept.md`](https://github.com/zwoo-hq/zwooc/blob/main/docs/concept.md)

### Run Mode

There a 3 available run modes: `run`, `watch`, `build`

### Projects

Define a sub-project with an adapter. The adapter will handle how commands are build. A project contains a number of profiles which can be run. The name of the project must equal the subpath.

#### Profiles

A profile is a run configuration for running a project in a certain run mode.

Profiles can be run via `zwoo <run|watch|build> <profile name>`

### Fragments

Fragments are individual commands that can be run before/with/after profiles. They are not bound to the adapter und are run with the folder of the project they are defined in. Fragments can adapt teh current run mode and profile.

Fragments can be run via `zwoo exec <fragment:configuration>`

### Compounds

Compounds are a combination of profiles that can be run together. They are defined in root scope.

Compounds can be started via `zwoo launch <compound name>`

### Genral Concepts

All runnable entities can define pre and post actions via `$pre` and `$post`


## Example

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
            "watch": "", // use a custom command
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
