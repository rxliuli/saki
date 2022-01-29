package globby

import (
	"github.com/rxliuli/saki/utils/array"
	"github.com/stretchr/testify/assert"
	"gopkg.in/ffmt.v1"
	"testing"
)

func TestSimpleGlob(t *testing.T) {
	matches := simpleGlob([]string{
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

func TestGlob(t *testing.T) {
	matches := Glob([]string{
		"apps/*",
		"apps/liuli-cli/templates/*",
		"libs/*",
		"archives/*",
		"examples/*",
		"scripts",
	}, Options{
		Cwd: "C:/Users/rxliuli/Code/Web/liuli-tools",
		Ignore: []string{
			"libs/async",
		},
	})
	_, _ = ffmt.Puts(matches)
	assert.NotEmpty(t, matches)
	assert.False(t, array.Contains(matches, "libs/async"))
}
