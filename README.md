# saki

> [中文](https://github.com/rxliuli/saki/blob/master/README.zh-CN.md)

I wanted to know how much performance can be improved by writing cli based on golang, so I tried to write this cli
application in golang.

## saki build

Build lib or cli programs based on esbuild, which is very fast.

````sh
saki build lib # build lib
saki build cli # build cli
````

Performance Testing

The following two cli are based on the build implemented by esbuild, but the real runtime of the latter is very long,
because nodejs itself also takes time (and is very slow) to load the code.

````sh
$ time saki build lib
Build completed: 13.2977ms

real 0m0.215s
user 0m0.061s
sys 0m0.105s

$ time saki build cli
Build completed: 86.9462ms

real 0m0.296s
user 0m0.075s
sys 0m0.106s
````

````sh
$ time pnpm liuli-cli build lib
Build complete: 29ms

real 0m1.043s
user 0m0.045s
sys 0m0.091s

$ time pnpm liuli-cli build cli
Build complete: 95ms

real 0m1.076s
user 0m0.030s
sys 0m0.152s
````

## saki run

An alternative to `pnpm --filter run` designed to improve the efficiency of running commands in multiple threads.

````sh
saki run setup # run the setup command on all modules (if there is one)
saki run --filter libs/* setup # run the setup command in all modules matching libs/*
# use --filter array
saki run --filter libs/* --filter apps/* setup # or use , to split
saki run --filter libs/*,apps/* setup
````

Performance Testing

pnpm + liuli-cli + dts

````sh
$ time pnpm --filter .run setup

real 0m45.186s
user 0m0.015s
sys 0m0.060s
````

pnpm + liuli-cli

````sh
$ time pnpm --filter .run setup

real 0m15.026s
user 0m0.030s
sys 0m0.045s
````

saki + liuli-cli (you can see here that although they are all multi-threaded, golang is twice as fast)

````sh
$ time saki run setup

real 0m7.166s
user 0m0.000s
sys 0m0.000s
````

pnpm + saki

````sh
$ time pnpm --filter .run setup

real 0m1.186s
user 0m0.030s
sys 0m0.045s
````

saki (almost unimaginable with a js toolchain)

````sh
$ time saki run setup

real 0m0.592s
user 0m0.000s
sys 0m0.015s
````

> Use git bash --filter parameter under Windows, please wrap it with `""`