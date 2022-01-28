package __tests__

import (
	"github.com/rxliuli/saki/utils/fsExtra"
	"github.com/stretchr/testify/assert"
	"gopkg.in/ffmt.v1"
	"path/filepath"
	"testing"
)

var rootDir string

func init() {
	rootDir = fsExtra.FindParentDir(fsExtra.Dirname(), func(dir string) bool {
		return fsExtra.PathExists(filepath.Join(dir, "go.mod"))
	})
}

func TestGlob(t *testing.T) {
	matches, err := filepath.Glob(filepath.Join(rootDir, "**", "*.go"))
	assert.NoError(t, err)
	_, _ = ffmt.Puts(matches)
}
