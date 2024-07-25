# ui overhaul todos

## restructuring types of ui

Static:

Since all tasks are now based of trees, there is no need to differentiate between "inline"/short running or lon running tasks. 

Options:
- whether command output is piped to stdout or not
- whether piped std out is prefixed


Interactive:

Here there should be a differentiation between short running in-line tasks (via exec/build) and long running tasks (watch/run)

Inline:

Display an inline progress update 


full screen:

- the known full screen tabs layout
- progress of pre/post is shown inside the tabs ()


Global Options:

- Quite Mode (default: false) - don't output anything


TODO: TaskTreeRunner: if cancel errors this error should also be returned from Start()
TODO: tree progress view show aggregated status