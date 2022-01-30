package runner

import (
	"fmt"
	"github.com/rxliuli/saki/builder"
	"github.com/rxliuli/saki/utils/array"
	"github.com/rxliuli/saki/utils/fsExtra"
	"github.com/rxliuli/saki/utils/globby"
	"github.com/rxliuli/saki/utils/object"
	"gopkg.in/yaml.v2"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"sync"
)

type Type string

type Program struct {
	WorkspaceType Type
	Cwd           string
}

type Module struct {
	Name    string            // 模块名称
	Deps    []string          // 依赖
	Path    string            // 模块绝对路径
	Scripts map[string]string // 脚本
}

type TaskState int

const (
	TaskStateWait TaskState = iota
	TaskStateRunning
	TaskStateSuccess
	TaskStateFailed
)

type Task struct {
	Module Module
	State  TaskState
}

//计算模块的依赖
func calcModulesDep(modules []Module) []Module {
	var nameSet = make(map[string]bool)
	for _, module := range modules {
		nameSet[module.Name] = true
	}
	var result []Module
	for _, module := range modules {
		resDeps := make([]string, 0)
		for _, name := range module.Deps {
			if nameSet[name] == true {
				resDeps = append(resDeps, name)
			}
		}
		result = append(result, Module{
			Name:    module.Name,
			Deps:    resDeps,
			Path:    module.Path,
			Scripts: module.Scripts,
		})
	}
	return result
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

	patterns := array.StringFilter(modules.Packages, func(s string) bool {
		return !strings.HasPrefix(s, "!")
	})
	ignores := append(array.StringMap(array.StringFilter(modules.Packages, func(s string) bool {
		return strings.HasPrefix(s, "!")
	}), func(str string) string {
		return strings.TrimPrefix(str, "!")
	}), []string{"node_modules", "**/node_modules", ".git", ".github", ".idea", ".*"}...)
	modulePaths := globby.Glob(patterns, globby.Options{
		Cwd:    receiver.Cwd,
		Ignore: ignores,
	})

	var res []Module
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
		res = append(res, Module{
			Name:    json.Name,
			Path:    absPath,
			Deps:    append(append(object.Keys(json.DevDependencies), object.Keys(json.Dependencies)...), object.Keys(json.PeerDependencies)...),
			Scripts: json.Scripts,
		})
	}
	return calcModulesDep(res)
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

func execTask(task Task, script string) TaskState {
	var command *exec.Cmd
	cmd := task.Module.Scripts[script]
	if cmd == "saki build lib" || cmd == "saki build cli" {
		command = exec.Command("saki", strings.Split(cmd, " ")[1:]...)
	} else {
		command = exec.Command("pnpm", "run", script)
	}
	command.Dir = task.Module.Path
	_, err := command.Output()
	if err != nil {
		_ = fmt.Errorf("[%s] 执行失败\n", task.Module.Name)
		return TaskStateFailed
	}
	fmt.Printf("[%s] 执行成功\n", task.Module.Name)
	return TaskStateSuccess
}

//计算模块的依赖
func removeTaskModuleDep(tasks []Task, complete string) {
	for i := range tasks {
		resDeps := make([]string, 0)
		for _, name := range tasks[i].Module.Deps {
			if name != complete {
				resDeps = append(resDeps, name)
			}
		}
		tasks[i].Module.Deps = resDeps
	}
}

func execTasks(tasks []Task, script string) {
	flag := true
	wg := sync.WaitGroup{}
	for i := range tasks {
		if tasks[i].State == TaskStateWait && len(tasks[i].Module.Deps) == 0 {
			tasks[i].State = TaskStateRunning
			flag = false
			wg.Add(1)
			go func(i int) {
				defer wg.Done()
				tasks[i].State = execTask(tasks[i], script)
				removeTaskModuleDep(tasks, tasks[i].Module.Name)
			}(i)
		}
	}
	wg.Wait()
	if flag {
		return
	}
	execTasks(tasks, script)
}

func (receiver Program) filterModuleByFilter(modules []Module, filters []string) []Module {
	if filters == nil || len(filters) == 0 {
		return modules
	}
	filterGlobs := globby.Patterns2Globs(array.StringMap(filters, func(str string) string {
		if filepath.IsAbs(str) {
			return str
		}
		return strings.ReplaceAll(filepath.Join(receiver.Cwd, str), "\\", "/")
	}))
	var res = make([]Module, 0)
	for _, module := range modules {
		if globby.MiniMatch(filterGlobs, strings.ReplaceAll(module.Path, "\\", "/")) || globby.MiniMatch(filterGlobs, module.Name) {
			res = append(res, module)
			continue
		}
	}
	return res
}

func (receiver Program) Run(options Options) {
	modules := calcModulesDep(receiver.filterModuleByFilter(filterModuleByScript(receiver.scanModules(), options.Script), options.Filter))
	fmt.Printf("扫描到 %d 个模块\n", len(modules))
	var tasks = make([]Task, len(modules))
	for i, module := range modules {
		tasks[i] = Task{
			Module: module,
			State:  TaskStateWait,
		}
	}
	execTasks(tasks, options.Script)
}
