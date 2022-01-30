package main

import (
	"errors"
	"github.com/rxliuli/saki/builder"
	"github.com/rxliuli/saki/runner"
	"github.com/rxliuli/saki/utils/array"
	"github.com/urfave/cli/v2"
	"gopkg.in/ffmt.v1"
	"log"
	"os"
	"strings"
)

func main() {
	cwd, _ := os.Getwd()

	program := builder.Program{
		Cwd: cwd,
	}
	app := &cli.App{
		Commands: cli.Commands{
			{
				Name:  "build",
				Usage: "构建命令",
				Subcommands: cli.Commands{
					{
						Name:  "lib",
						Usage: "构建 lib",
						Action: func(context *cli.Context) error {
							program.BuildLib()
							return nil
						},
					},
					{
						Name:  "cli",
						Usage: "构建 cli",
						Action: func(context *cli.Context) error {
							program.BuildCli()
							return nil
						},
					},
					{
						Name:  "single",
						Usage: "构建 cli",
						Flags: []cli.Flag{
							&cli.StringSliceFlag{
								Name:     "target",
								Usage:    "输出目标，可选值 esm/cjs/iife/cli",
								Required: true,
							},
						},
						Action: func(context *cli.Context) error {
							_, _ = ffmt.Puts("context: ", context.StringSlice("target"))
							return nil
						},
					},
				},
			},
			{
				Name:  "run",
				Usage: "运行命令",
				Flags: []cli.Flag{
					&cli.StringSliceFlag{
						Name:     "filter",
						Usage:    "模块过滤器路径列表，可使用 glob 模式",
						Required: false,
					},
				},
				Action: func(context *cli.Context) error {
					cmd := context.Args().First()
					if cmd == "" {
						return errors.New("请输入运行的命令")
					}
					filters := array.StringFlatMap(context.StringSlice("filter"), func(s string) []string {
						return strings.Split(s, ",")
					})
					//ffmt.Puts("filters: ", cmd, filters, context.StringSlice("filter"))
					runner.Program{
						Cwd: cwd,
					}.Run(runner.Options{
						Filter: filters,
						Script: cmd,
					})
					return nil
				},
			},
		},
		Name:    "saki",
		Usage:   "基于 esbuild 实现高层次的构建工具",
		Version: "0.1.0",
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}
