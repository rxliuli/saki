package globby

import (
	"github.com/stretchr/testify/assert"
	"gopkg.in/ffmt.v1"
	"testing"
)

func TestReal(t *testing.T) {
	matches := Glob([]string{
		"apps/*",
		"apps/liuli-cli/templates/*",
		"libs/*",
		"archives/*",
		"examples/*",
		"scripts",
	}, Options{
		Cwd:    "C:/Users/rxliuli/Code/Web/liuli-tools",
		Ignore: []string{},
	})
	_, _ = ffmt.Puts(matches)
	assert.NotEmpty(t, matches)
}
