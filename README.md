# Raft

This is a [Raft](https://raft.github.io/) protocol implementation in Go.

## Build

Build it using a regular `go build .` command.

## Usage

To see this implementation in action you need at least 2 (3 is better) nodes.

Start the first node by just supplying the listen address (if omitted defults to `localhost`) and port of the UDP server:

```
./raft -port 3000 [-addr 127.0.0.1]
```

Then add more nodes by specifying the the host and port of the first server:

```
./raft -port 3001 -phost 127.0.0.1 -pport 3000
```

Other nodes can be easily added to the cluster by specifying host and port of any running node. New nodes are propogated to the entire cluster to keep everyone in sync.

```
./raft -port 3002 -phost 127.0.0.1 -pport 3001
```