package main

import (
	"embed"
	"encoding/json"
	"errors"
	"github.com/rxliuli/saki/builder"
	"github.com/rxliuli/saki/runner"
	"github.com/rxliuli/saki/utils/array"
	"github.com/urfave/cli/v2"
	"log"
	"os"
	"strings"
)

//go:embed npm/package.json
var packageJson embed.FS

type PackageJson struct {
	Version string `json:"version"`
}

func readVersion() string {
	file, _ := packageJson.ReadFile("npm/package.json")
	var pkgJson PackageJson
	_ = json.Unmarshal(file, &pkgJson)
	return pkgJson.Version
}

func main() {
	version := readVersion()
	cwd, _ := os.Getwd()

	program := builder.Program{
		Cwd: cwd,
	}
	app := &cli.App{
		Commands: cli.Commands{
			{
				Name:  "build",
				Usage: "构建命令",
				Flags: []cli.Flag{
					&cli.BoolFlag{
						Name:    "watch",
						Usage:   "监听文件变化，自动编译",
						Value:   false,
						Aliases: []string{"w"},
					},
				},
				Subcommands: cli.Commands{
					{
						Name:  "lib",
						Usage: "构建 lib",
						Action: func(context *cli.Context) error {
							program.Watch = context.Bool("watch")
							return program.BuildToTargets([]builder.Target{builder.TargetEsm, builder.TargetCjs})
						},
					},
					{
						Name:  "cli",
						Usage: "构建 cli",
						Action: func(context *cli.Context) error {
							program.Watch = context.Bool("watch")
							return program.BuildToTargets([]builder.Target{builder.TargetCli, builder.TargetEsm, builder.TargetCjs})
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
							program.Watch = context.Bool("watch")
							return program.BuildToTargets(array.StringFlatMap(context.StringSlice("target"), func(s string) []string {
								return strings.Split(s, ",")
							}))
						},
					},
				},
			},
			{
				Name:   "run",
				Usage:  "运行命令",
				Hidden: true,
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
		Version: version,
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}
