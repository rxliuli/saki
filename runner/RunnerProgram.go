package runner

import (
	"fmt"
	"github.com/rxliuli/saki/builder"
	"github.com/rxliuli/saki/utils/fsExtra"
	"github.com/rxliuli/saki/utils/globby"
	"github.com/rxliuli/saki/utils/object"
	"gopkg.in/yaml.v2"
	"os"
	"os/exec"
	"path/filepath"
)

type Type string

type Program struct {
	WorkspaceType Type
	Cwd           string
}

type Module struct {
	Name    string
	Deps    []string
	Path    string
	Scripts map[string]string
}

type TaskState int

const (
	TaskStateWait TaskState = iota
	TaskStateSuccess
	TaskStateFailed
)

type Task struct {
	Module Module
	State  TaskState
}

func (receiver Program) scanModules() []Module {
	file, err := os.ReadFile(filepath.Join(receiver.Cwd, "pnpm-workspace.yaml"))
	if err != nil {
		panic("扫描所有模块失败")
	}
	var modules struct {
		Packages []string `json:"packages"`
	}
	err = yaml.Unmarshal(file, &modules)
	if err != nil {
		panic("解析 yaml 文件失败")
	}
	modulePaths := globby.Glob(modules.Packages, globby.Options{
		Cwd:    receiver.Cwd,
		Ignore: []string{},
	})

	var res []Module
	var nameSet = make(map[string]bool)
	for _, path := range modulePaths {
		absPath := filepath.Join(receiver.Cwd, path)
		var json builder.PackageJson
		jsonPath := filepath.Join(absPath, "package.json")
		if !fsExtra.PathExists(jsonPath) {
			continue
		}
		err := fsExtra.ReadJson(jsonPath, &json)
		if err != nil {
			panic(fmt.Sprintf("解析 json 失败 %s", absPath))
		}
		nameSet[json.Name] = true
		res = append(res, Module{
			Name:    json.Name,
			Path:    absPath,
			Deps:    append(append(object.Keys(json.DevDependencies), object.Keys(json.Dependencies)...), object.Keys(json.PeerDependencies)...),
			Scripts: json.Scripts,
		})
	}
	for i, module := range res {
		var resDeps []string
		for _, name := range module.Deps {
			if nameSet[name] == true {
				resDeps = append(resDeps, name)
			}
		}
		module.Deps = resDeps
		res[i] = module
	}
	return res
}

type Options struct {
	Filter []string
	Script string
}

func filterModuleByScript(modules []Module, script string) []Module {
	var res []Module
	for _, module := range modules {
		if module.Scripts[script] != "" {
			res = append(res, module)
		}
	}
	return res
}

func execTasks(tasks []Task, script string) {
	for i, task := range tasks {
		if task.State == TaskStateWait && len(task.Module.Deps) == 0 {
			err := exec.Command("pnpm run " + script).Run()
			if err != nil {
				_ = fmt.Errorf("[%s] 执行失败", task.Module.Name)
				task.State = TaskStateFailed
				continue
			}
			fmt.Printf("[%s] 执行成功", task.Module.Name)
			task.State = TaskStateSuccess
			tasks[i] = task
		}
	}
	execTasks(tasks, script)
}

func (receiver Program) Run(options Options) {
	modules := filterModuleByScript(receiver.scanModules(), options.Script)
	var tasks = make([]Task, len(modules))
	for i, module := range modules {
		tasks[i] = Task{
			Module: module,
			State:  TaskStateWait,
		}
	}
	execTasks(tasks, options.Script)
}
