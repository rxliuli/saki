# saki

想知道基于 golang 编写 cli 能够提高多少性能，所以尝试使用 golang 编写了这个 cli 应用。

## saki build

基于 esbuild 构建 lib 或 cli 程序，它非常快。

```sh
saki build lib # 构建 lib
saki build cli # 构建 cli
```

性能测试

下面两个 cli 都基于 esbuild 实现的构建，但后者真实运行时间很长，因为 nodejs 本身加载代码也需要时间（而且很慢）。

```sh
$ time saki build lib
构建完成: 13.2977ms

real    0m0.215s
user    0m0.061s
sys     0m0.105s

$ time saki build cli
构建完成: 86.9462ms

real    0m0.296s
user    0m0.075s
sys     0m0.106s
```

```sh
$ time pnpm liuli-cli build lib
构建完成: 29ms

real    0m1.043s
user    0m0.045s
sys     0m0.091s

$ time pnpm liuli-cli build cli
构建完成: 95ms

real    0m1.076s
user    0m0.030s
sys     0m0.152s
```

## saki run

一个 `pnpm --filter run` 的替代命令，旨在提高多线程运行命令的效率。

```sh
saki run setup # 在所有模块运行 setup 命令（如果有这个命令）
saki run --filter libs/* setup # 在所有匹配 libs/* 的模块中运行 setup 命令
# 使用 --filter 数组
saki run --filter libs/* --filter apps/* setup # 或者使用 , 分割
saki run --filter libs/*,apps/* setup
```

性能测试

pnpm + cli + dts

```sh
$ time pnpm --filter . run setup

real    0m45.186s
user    0m0.015s
sys     0m0.060s
```

pnpm + cli

```sh
$ time pnpm --filter . run setup

real    0m15.026s
user    0m0.030s
sys     0m0.045s
```

saki + cli（这里可以看到虽然都是多线程，但 golang 的快一倍）

```sh
$ time saki run setup

real    0m7.166s
user    0m0.000s
sys     0m0.000s
```

pnpm + saki

```sh
$ time pnpm --filter . run setup

real    0m1.186s
user    0m0.030s
sys     0m0.045s
```

saki（几乎是使用 js 工具链难以想象的）

```sh
$ time saki run setup

real    0m0.592s
user    0m0.000s
sys     0m0.015s
```

> Windows 下使用 git bash --filter 参数请用 `""` 包裹
