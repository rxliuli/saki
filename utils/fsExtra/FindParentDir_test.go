package fsExtra

import (
	"github.com/stretchr/testify/assert"
	"path/filepath"
	"testing"
)

func TestFindParentDir(t *testing.T) {
	findParentDir := FindParentDir(Dirname(), func(dir string) bool {
		return PathExists(filepath.Join(dir, "go.mod"))
	})
	assert.NotEqual(t, findParentDir, "")
	assert.True(t, PathExists(filepath.Join(findParentDir, "pnpm-lock.yaml")))
}
