package __tests__

import (
	"github.com/rxliuli/saki/utils/fsExtra"
	"github.com/yargevad/filepathx"
	"gopkg.in/ffmt.v1"
	"path"
	"path/filepath"
	"strings"
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
	if err != nil {
		t.Error(err)
	}
	_, _ = ffmt.Puts(matches)
}

func TestFilePathX(t *testing.T) {
	matches, err := filepathx.Glob(filepath.Join(rootDir, "**/*.go"))
	if err != nil {
		t.Error(err)
	}
	_, _ = ffmt.Puts(matches)
	for _, match := range matches {
		if !strings.HasSuffix(match, ".go") {
			t.Error(match)
		}
	}
}

func TestFilePathXByRealProject(t *testing.T) {
	matches, err := filepathx.Glob(filepath.Join("C:/Users/rxliuli/Code/Pkg/liuli-tools", "libs/*", "package.json"))
	if err != nil {
		t.Error(err)
	}
	//_, _ = ffmt.Puts("matches:", matches)
	var res []string
	for _, match := range matches {
		isMatch, err := path.Match("libs/async", match)
		if err != nil {
			return
		}
		if !isMatch {
			res = append(res, match)
		}
	}
	ffmt.Puts("res:", res)
}
