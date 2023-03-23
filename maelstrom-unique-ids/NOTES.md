# Unique Id Generator

When running a unique id generator on single instance it works fine. The problem is this approach is not scalable
On a single node, we can use atomic increment and every time we as asked for id we inc and return the value
Even this approach has limit since at max we can vend 2^64 - 1

docker run --name pg-dist-sys -p 5432:5432 -e POSTGRES_PASSWORD=password -d postgres:14