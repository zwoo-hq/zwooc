# runner

The task runner used in `zwooc` runs tasks based on a tree like structure. As displayed in in `zwooc graph` or using the `--dry-run` option, zwooc resolves dependencies between tasks into a tree like structure. 

Each tree node consists of an task that gets executed and a number of preceding and succeeding tasks. In this model each sub tree own its acts as an pre or post dependency asa a whole. Executing those subtrees independently allows for most optimal scheduling without unnecessary extra waiting times caused by cross dependencies.

This new scheduling model comes with a few difficulties when displayed in a TUI since the sequence itself cant be divided into separate stages. Thus zwooc is not able to display any exact progress information at the moment.