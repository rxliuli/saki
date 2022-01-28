package glob

import (
	"github.com/rxliuli/saki/utils/fsExtra"
	"github.com/stretchr/testify/assert"
	"gopkg.in/ffmt.v1"
	"path/filepath"
	"testing"
)

var rootPath string

func init() {
	rootPath = fsExtra.FindParentDir(fsExtra.Dirname(), func(dir string) bool {
		return fsExtra.PathExists(filepath.Join(dir, "go.mod"))
	})
}

func TestMatch(t *testing.T) {
	assert.True(t, match([]string{".*"}, ".env"))
	assert.False(t, match([]string{".*"}, "test.go"))
	assert.False(t, match([]string{".*", "node_modules"}, "test.go"))
	assert.True(t, match([]string{".*", "node_modules"}, "node_modules"))
	assert.True(t, match([]string{"libs/*"}, "libs/array"))
	assert.False(t, match([]string{"libs/*"}, "libs/array/src"))
}

func TestMatchByPath(t *testing.T) {
	assert.True(t, match([]string{"*.go"}, "utils/array/array.go"))
	assert.True(t, match([]string{"*_test.go"}, "utils/array/array_test.go"))
	assert.False(t, match([]string{"*_test.go"}, "utils/array/array.go"))
}

func TestGlobByFile(t *testing.T) {
	matches, err := Glob([]string{"**/*.go"}, Options{
		Cwd:    rootPath,
		Ignore: []string{"**/*_test.go", "node_modules", ".*"},
	})
	assert.NoError(t, err)
	_, _ = ffmt.Puts(matches)
}

func TestGlobByDir(t *testing.T) {
	matches, err := Glob([]string{"**/__test__"}, Options{
		Cwd:    rootPath,
		Ignore: []string{"**/*_test.go", "node_modules", ".*"},
	})
	assert.NoError(t, err)
	_, _ = ffmt.Puts(matches)
}
