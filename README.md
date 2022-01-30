# saki

想知道基于 golang 编写 cli 能够提高多少性能，所以尝试使用 golang 编写了这个 cli 应用。

## saki build

```sh
saki build lib # 构建 lib
saki build cli # 构建 cli
```

## saki run

```sh
saki run setup # 在所有模块运行 setup 命令（如果有这个命令）
saki run --filter libs/* setup # 在所有匹配 libs/* 的模块中运行 setup 命令
# 使用 --filter 数组
saki run --filter libs/* --filter apps/* setup # 或者使用 , 分割
saki run --filter libs/*,apps/* setup
```

> Windows 下使用 git bash --filter 参数请用 `""` 包裹
