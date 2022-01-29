package runner

import (
	"github.com/stretchr/testify/assert"
	"gopkg.in/ffmt.v1"
	"path/filepath"
	"testing"
)

var program Program

func init() {
	program = Program{
		Cwd: "C:/Users/rxliuli/Code/Web/liuli-tools",
	}
}

func TestScanModules(t *testing.T) {
	modules := program.scanModules()
	_, _ = ffmt.Puts(modules)
	assert.NotEmpty(t, modules)
}

func TestCalcModulesDep(t *testing.T) {
	res := calcModulesDep([]Module{
		{Name: "a", Deps: []string{"b"}},
		{Name: "b", Deps: []string{}},
		{Name: "c", Deps: []string{"d"}},
	})
	exp := []Module{
		{Name: "a", Deps: []string{"b"}},
		{Name: "b", Deps: []string{}},
		{Name: "c", Deps: []string{}},
	}
	assert.Equal(t, res, exp)
}

func TestRemoveTaskModuleDep(t *testing.T) {
	tasks := []Task{
		{Module: Module{Name: "a", Deps: []string{"b"}}},
		{Module: Module{Name: "b", Deps: []string{}}},
	}
	removeTaskModuleDep(tasks, "b")
	assert.Equal(t, tasks[0].Module.Deps, []string{})
}

func TestExecTask(t *testing.T) {
	execTask(Task{
		Module: Module{
			Name: "@liuli-tools/async",
			Scripts: map[string]string{
				"setup": "saki build lib",
			},
			Path: filepath.Join(program.Cwd, "libs/async"),
			Deps: []string{},
		},
	}, "setup")

}

func TestExecTasks(t *testing.T) {
	execTasks([]Task{
		{
			State: TaskStateWait,
			Module: Module{
				Name: "@liuli-tools/async",
				Scripts: map[string]string{
					"setup": "saki build lib",
				},
				Path: filepath.Join(program.Cwd, "libs/async"),
				Deps: []string{},
			},
		},
		{
			State: TaskStateWait,
			Module: Module{
				Name: "@liuli-tools/array",
				Scripts: map[string]string{
					"setup": "saki build lib",
				},
				Path: filepath.Join(program.Cwd, "libs/array"),
				Deps: []string{},
			},
		},
		{
			State: TaskStateWait,
			Module: Module{
				Name: "@liuli-tools/i18next-dts-gen",
				Scripts: map[string]string{
					"setup": "saki build cli",
				},
				Path: filepath.Join(program.Cwd, "libs/i18next-dts-gen"),
				Deps: []string{},
			},
		},
		{
			State: TaskStateWait,
			Module: Module{
				Name: "@liuli-util/cli",
				Scripts: map[string]string{
					"setup": "saki build cli",
				},
				Path: filepath.Join(program.Cwd, "apps/liuli-cli"),
				Deps: []string{},
			},
		},
	}, "setup")
}

func TestProgram_Run(t *testing.T) {
	program.Run(Options{
		Filter: []string{"./libs/*"},
		Script: "setup",
	})
}
