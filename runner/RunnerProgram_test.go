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
