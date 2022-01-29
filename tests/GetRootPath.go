package tests

import (
	"github.com/rxliuli/saki/utils/fsExtra"
	"path/filepath"
)

func GetRootPath() string {
	return fsExtra.FindParentDir(fsExtra.Dirname(), func(dir string) bool {
		return fsExtra.PathExists(filepath.Join(dir, "go.mod"))
	})
}
