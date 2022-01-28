package runner

import (
	"fmt"
	"github.com/rxliuli/saki/builder"
	"github.com/rxliuli/saki/utils/fsExtra"
	"github.com/rxliuli/saki/utils/globby"
	"github.com/rxliuli/saki/utils/object"
	"gopkg.in/yaml.v2"
	"os"
	"path/filepath"
)

type Type string

type Program struct {
	WorkspaceType Type
	Cwd           string
}

type Module struct {
	Name string
	Deps []string
	Path string
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
			Name: json.Name,
			Path: absPath,
			Deps: append(append(object.Keys(json.DevDependencies), object.Keys(json.Dependencies)...), object.Keys(json.PeerDependencies)...),
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
