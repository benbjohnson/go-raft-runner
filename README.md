Raft Runner
==============

## Overview

The Raft runner is a testbed for running the [go-raft](https://github.com/benbjohnson/go-raft) library under different scenarios.


## Running

To install:

```sh
$ go get github.com/benbjohnson/go-raft-runner
```

The first instance of the runner can be started with no options:

```sh
$ go-raft-runner
Running on localhost:20000
```

Subsequent instances of the runner should point to a node in the cluster:

```sh
$ go-raft-runner localhost:20000
Running on localhost:20001
```

