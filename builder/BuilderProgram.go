package builder

import (
	"errors"
	"fmt"
	"github.com/evanw/esbuild/pkg/api"
	"github.com/rxliuli/saki/builder/plugin"
	"github.com/rxliuli/saki/utils/fsExtra"
	"github.com/rxliuli/saki/utils/object"
	"github.com/swaggest/assertjson/json5"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"
)

type Program struct {
	Cwd   string
	Watch bool
}

var (
	globalExternal = []string{"esbuild", "pnpapi", "ts-morph", "ssh2"}
)

type PackageJson struct {
	Name             string            `json:"name"`
	Dependencies     map[string]string `json:"dependencies"`
	DevDependencies  map[string]string `json:"devDependencies"`
	PeerDependencies map[string]string `json:"peerDependencies"`
	Scripts          map[string]string `json:"scripts"`
}

func (receiver Program) getDeps() []string {
	json := PackageJson{}
	err := fsExtra.ReadJson(filepath.Join(receiver.Cwd, "package.json"), &json)
	if err != nil {
		panic("解析 package.json 失败")
	}
	var res []string
	if json.Dependencies != nil {
		res = append(res, object.Keys(json.Dependencies)...)
	}
	if json.DevDependencies != nil {
		res = append(res, object.Keys(json.DevDependencies)...)
	}
	if json.PeerDependencies != nil {
		res = append(res, object.Keys(json.PeerDependencies)...)
	}
	return res
}

type TsConfigCompilerOptions struct {
	Lib []string `json:"lib"`
}

type TsConfig struct {
	CompilerOptions TsConfigCompilerOptions `json:"compilerOptions"`
}

func (receiver Program) getPlatform() api.Platform {
	tsconfigPath := filepath.Join(receiver.Cwd, "tsconfig.json")
	if fsExtra.PathExists(tsconfigPath) {
		var tsconfigJson TsConfig
		file, err := os.ReadFile(tsconfigPath)
		if err != nil {
			panic("无法读取文件 tsconfig.json")
		}
		err = json5.Unmarshal(file, &tsconfigJson)
		if err != nil {
			panic("无法解析 tsconfig.json")
		}
		for _, lib := range tsconfigJson.CompilerOptions.Lib {
			if strings.ToLower(lib) == "dom" {
				return api.PlatformBrowser
			}
		}
	}
	pkgPath := filepath.Join(receiver.Cwd, "package.json")
	if fsExtra.PathExists(pkgPath) {
		var pkgJson PackageJson
		err := fsExtra.ReadJson(pkgPath, &pkgJson)
		if err != nil {
			panic("无法解析 package.json")
		}
		for k := range pkgJson.DevDependencies {
			if k == "@types/node" {
				return api.PlatformNode
			}
		}
	}
	return api.PlatformNeutral
}

func (receiver Program) getPlugins(platform api.Platform) []api.Plugin {
	var plugins []api.Plugin
	if platform == api.PlatformBrowser {
	}
	if platform == api.PlatformNode {
		plugins = append(plugins, plugin.NodeExternal())
	}
	if receiver.Watch {
		plugins = append(plugins, plugin.Logger())
	}
	return plugins
}

func (receiver Program) getBaseOptions() api.BuildOptions {
	var external = globalExternal
	if receiver.Watch {
		external = append(external, receiver.getDeps()...)
	}
	var watch api.WatchMode
	if receiver.Watch {
		watch = api.WatchMode{}
	}
	platform := receiver.getPlatform()
	return api.BuildOptions{
		AbsWorkingDir:     receiver.Cwd,
		Sourcemap:         api.SourceMapExternal,
		Bundle:            true,
		Watch:             &watch,
		External:          globalExternal,
		Platform:          platform,
		MinifyWhitespace:  !receiver.Watch,
		MinifyIdentifiers: !receiver.Watch,
		MinifySyntax:      !receiver.Watch,
		Incremental:       receiver.Watch,
		Metafile:          receiver.Watch,
		Write:             true,
		Plugins:           receiver.getPlugins(platform),
	}
}

func (receiver Program) getEsmOptions() api.BuildOptions {
	options := receiver.getBaseOptions()
	options.EntryPoints = []string{filepath.Join(receiver.Cwd, "src/index.ts")}
	options.Outfile = filepath.Join(receiver.Cwd, "dist/index.esm.js")
	options.Format = api.FormatESModule
	options.Plugins = append(options.Plugins, plugin.AutoExternal())
	return options
}

func (receiver Program) getCjsOptions() api.BuildOptions {
	options := receiver.getBaseOptions()
	options.EntryPoints = []string{filepath.Join(receiver.Cwd, "src/index.ts")}
	options.Outfile = filepath.Join(receiver.Cwd, "dist/index.js")
	options.Format = api.FormatCommonJS
	options.Plugins = append(options.Plugins, plugin.AutoExternal())
	return options
}

func (receiver Program) getIifeOptions() api.BuildOptions {
	options := receiver.getBaseOptions()
	options.EntryPoints = []string{filepath.Join(receiver.Cwd, "src/index.ts")}
	options.Outfile = filepath.Join(receiver.Cwd, "dist/index.iife.js")
	options.Format = api.FormatIIFE
	return options
}

func (receiver Program) getCliOptions() api.BuildOptions {
	options := receiver.getBaseOptions()
	options.EntryPoints = []string{filepath.Join(receiver.Cwd, "src/bin.ts")}
	options.Outfile = filepath.Join(receiver.Cwd, "dist/bin.js")
	options.Format = api.FormatCommonJS
	options.Platform = api.PlatformNode
	options.Banner = map[string]string{
		"js": "#!/usr/bin/env node",
	}
	return options
}

type Target = string

const (
	TargetEsm  Target = "esm"
	TargetCjs  Target = "cjs"
	TargetIife Target = "iife"
	TargetCli  Target = "cli"
)

func (receiver Program) GetOptionsByTarget(target Target) api.BuildOptions {
	switch target {
	case TargetEsm:
		return receiver.getEsmOptions()
	case TargetCjs:
		return receiver.getCjsOptions()
	case TargetIife:
		return receiver.getIifeOptions()
	case TargetCli:
		return receiver.getCliOptions()
	default:
		panic("无法识别的目标")
	}
}

func resolveResultError(results []api.BuildResult) error {
	for _, result := range results {
		if len(result.Errors) != 0 {
			message := result.Errors[0]
			location := message.Location
			return fmt.Errorf("构建失败: %s %s:%d:%d\n", message.Text, location.File, location.Line, location.Column)
		}
	}
	return nil
}

func (receiver Program) BuildToTargets(targets []Target) error {
	start := time.Now()
	if !receiver.Watch {
		err := os.RemoveAll(filepath.Join(receiver.Cwd, "dist"))
		if err != nil {
			_ = fmt.Errorf("无法清理 dist 目录")
		}
	}
	var res = make([]api.BuildResult, len(targets))
	var wg sync.WaitGroup
	for i, target := range targets {
		wg.Add(1)
		go func(i int, target string) {
			defer wg.Done()
			res[i] = api.Build(receiver.GetOptionsByTarget(target))
		}(i, target)
	}
	wg.Wait()
	if receiver.Watch {
		<-make(chan bool)
	} else {
		fmt.Printf("构建完成: %s\n", time.Now().Sub(start).String())
	}
	return resolveResultError(res)
}

func (receiver Program) absPath(path string) string {
	if filepath.IsAbs(path) {
		return path
	}
	return filepath.Join(receiver.Cwd, path)
}

type SourceMap string

//goland:noinspection GoUnusedConst
const (
	SourceMapNone              SourceMap = "none"
	SourceMapInline            SourceMap = "inline"
	SourceMapLinked            SourceMap = "linked"
	SourceMapExternal          SourceMap = "external"
	SourceMapInlineAndExternal SourceMap = "inlineAndExternal"
)

type Platform string

//goland:noinspection GoUnusedConst
const (
	PlatformBrowser Platform = "browser"
	PlatformNode    Platform = "node"
	PlatformNeutral Platform = "neutral"
)

type Format string

//goland:noinspection GoUnusedConst
const (
	FormatDefault  Format = "default"
	FormatIIFE     Format = "iife"
	FormatCommonJS Format = "cjs"
	FormatESModule Format = "esm"
)

type PluginName string

const (
	PluginAutoExternal PluginName = "auto-external" // 自动排除所有依赖
	PluginLogger       PluginName = "logger"        // 打印构建日志
	PluginEnv          PluginName = "env"           // 替换环境变量，支持 import.meta.env 和 process.env
	PluginNodeExternal PluginName = "node-external" // 排除 node:* 原生模块
	PluginRaw          PluginName = "raw"           // 支持引入文本文件
)

// BuildOptions 支持 esbuild 构建选项的子集
type BuildOptions struct {
	EntryPoints []string                   `json:"entryPoints"`
	Outfile     string                     `json:"outfile"`
	Platform    Platform                   `json:"platform"`
	Format      Format                     `json:"format"`
	Sourcemap   SourceMap                  `json:"sourcemap"`
	Minify      bool                       `json:"minify"`
	External    []string                   `json:"external"`
	Bundle      bool                       `json:"bundle"`
	Plugins     map[PluginName]interface{} `json:"plugins"` // 插件配置，插件名 => 配置的 Map
}

func convertSourceMap(sourceMap SourceMap) api.SourceMap {
	switch sourceMap {
	case SourceMapInline:
		return api.SourceMapInline
	case SourceMapLinked:
		return api.SourceMapLinked
	case SourceMapExternal:
		return api.SourceMapExternal
	case SourceMapInlineAndExternal:
		return api.SourceMapInlineAndExternal
	default:
		return api.SourceMapNone
	}
}

func convertPlatform(platform Platform) api.Platform {
	switch platform {
	case PlatformBrowser:
		return api.PlatformBrowser
	case PlatformNode:
		return api.PlatformNode
	default:
		return api.PlatformNeutral
	}
}

func convertFormat(format Format) api.Format {
	switch format {
	case FormatIIFE:
		return api.FormatIIFE
	case FormatCommonJS:
		return api.FormatCommonJS
	case FormatESModule:
		return api.FormatESModule
	default:
		return api.FormatDefault
	}
}

func convertPlugin(name PluginName) api.Plugin {
	switch name {
	case PluginAutoExternal:
		return plugin.AutoExternal()
	case PluginLogger:
		return plugin.Logger()
	case PluginEnv:
		return plugin.Env()
	case PluginNodeExternal:
		return plugin.NodeExternal()
	case PluginRaw:
		return plugin.Raw()
	default:
		panic("无法识别的插件")
	}
}

func convertPlugins(plugins map[PluginName]interface{}) []api.Plugin {
	res := make([]api.Plugin, 0)
	for k := range plugins {
		res = append(res, convertPlugin(k))
	}
	return res
}

//计算 esbuild 构建选项
func (receiver Program) calcOptions(options BuildOptions) api.BuildOptions {
	for i := range options.EntryPoints {
		options.EntryPoints[i] = receiver.absPath(options.EntryPoints[i])
	}
	if options.Outfile != "" {
		options.Outfile = receiver.absPath(options.Outfile)
	}
	var watch api.WatchMode
	if receiver.Watch {
		watch = api.WatchMode{}
	}
	baseOptions := api.BuildOptions{
		Write:         true,
		Bundle:        true,
		AbsWorkingDir: receiver.Cwd,
		Watch:         &watch,
	}
	if receiver.Watch {
		baseOptions.Plugins = append(baseOptions.Plugins, plugin.Logger())
	}
	baseOptions.EntryPoints = options.EntryPoints
	baseOptions.Outfile = options.Outfile
	baseOptions.Platform = convertPlatform(options.Platform)
	baseOptions.Format = convertFormat(options.Format)
	baseOptions.Sourcemap = convertSourceMap(options.Sourcemap)
	if options.Minify {
		baseOptions.MinifyWhitespace = true
		baseOptions.MinifySyntax = true
		baseOptions.MinifyIdentifiers = true
	}
	baseOptions.External = options.External
	baseOptions.Bundle = options.Bundle
	baseOptions.Plugins = append(baseOptions.Plugins, convertPlugins(options.Plugins)...)
	return baseOptions
}

func loadConfig(cwd string) (error, []BuildOptions) {
	var json struct {
		Saki []BuildOptions `json:"saki"`
	}
	err := fsExtra.ReadJson(filepath.Join(cwd, "package.json"), &json)
	if err != nil {
		return errors.New("解析 package.json 失败"), nil
	}
	if json.Saki == nil {
		return errors.New("没有找到 saki 字段配置项"), nil
	}
	return nil, json.Saki
}

func (receiver Program) buildByConfig(options []BuildOptions) error {
	start := time.Now()
	for _, options := range options {
		result := api.Build(receiver.calcOptions(options))
		if len(result.Errors) != 0 {
			return errors.New(result.Errors[0].Text)
		}
	}
	if receiver.Watch {
		<-make(chan bool)
	} else {
		fmt.Printf("构建完成: %s\n", time.Now().Sub(start).String())
	}
	return nil
}

func (receiver Program) BuildByConfig() error {
	err, configs := loadConfig(receiver.Cwd)
	if err != nil {
		return err
	}
	return receiver.buildByConfig(configs)
}
