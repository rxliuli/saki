package runner

import (
	"github.com/stretchr/testify/assert"
	"gopkg.in/ffmt.v1"
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

func TestExecTask(t *testing.T) {
	execTasks([]Task{
		{
			State: TaskStateWait,
			Module: Module{
				Name: "",
				Scripts: map[string]string{
					"setup": "go build",
				},
				Deps: []string{},
			},
		},
	}, "setup")
}
