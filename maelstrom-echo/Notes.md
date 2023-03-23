./maelstrom test -w echo --bin ~/go/bin/maelstrom-echo --node-count 1 --time-limit 10
This command instructs maelstrom to run the "echo" workload against our binary.
It runs a single node and it will send "echo" commands for 10 seconds.

Note:

Maelstrom will only inject network failures and it will not intentionally crash your node process
so you don't need to worry about persistence.
You can use in-memory data structures for these challenges.
