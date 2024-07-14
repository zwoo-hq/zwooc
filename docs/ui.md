# ui

`zwooc` comes with two general flavors of ui: an `interactive` one and a `non-interactive` one.

The non interactive ui is best suited for CI and non TTY environments. It keeps the output in a clean and readonly log format to enable easy debugging of build failures. While possible, its not perfectly suited for watch mode.

The interactive ui is best for development. It communicates the state and progress of tasks clearly in real time while also providing an efficient way to access standard out of running tasks. On top of that, the interactive mode allows interactions such as restarting, scheduling or stopping of tasks.