package main

import (
	"github.com/urfave/cli/v2"
	"gopkg.in/ffmt.v1"
	"os"
	"saki/build"
)

func main() {
	cwd, _ := os.Getwd()

	program := build.BuilderProgram{
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
		},
		Name:    "saki",
		Usage:   "基于 esbuild 实现高层次的构建工具",
		Version: "0.1.0",

		Action: func(context *cli.Context) error {
			_, _ = ffmt.Puts("context: ", context)
			return nil
		},
	}

	_ = app.Run(os.Args)
}
