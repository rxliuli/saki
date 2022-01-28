package fsExtra

import (
	"path/filepath"
	"testing"
)

func TestFindParentDir(t *testing.T) {
	findParentDir := FindParentDir(Dirname(), func(dir string) bool {
		return PathExists(filepath.Join(dir, "go.mod"))
	})
	if findParentDir == "" {
		t.Error("FindParentDir failed")
	}
	if !PathExists(filepath.Join(findParentDir, "pnpm-lock.yaml")) {
		t.Error("FindParentDir failed")
	}
}
