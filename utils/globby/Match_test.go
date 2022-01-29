package globby

import (
	"github.com/gobwas/glob"
	"github.com/rxliuli/saki/tests"
	"github.com/rxliuli/saki/utils/array"
	"github.com/stretchr/testify/assert"
	"gopkg.in/ffmt.v1"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func match(p string, s string) bool {
	return glob.MustCompile(p, '/').Match(s)
}

func TestGlobMatch(t *testing.T) {
	assert.True(t, match("libs/*", "libs/test"))
	assert.False(t, match("libs/*", "libs/test/test"))
	assert.True(t, match("libs/**", "libs/test/test"))
	assert.True(t, match("libs/*/package.json", "libs/async/package.json"))
	assert.True(t, match("libs/**/package.json", "libs/async/package.json"))
	assert.True(t, match(".*", ".git"))
	assert.False(t, match(".*", "src"))
	assert.True(t, match("*.go", "main.go"))
	assert.True(t, match("**/*.go", "src/util/main.go"))
	assert.False(t, match("*.go", "src/main.go"))
	assert.True(t, match("**/node_modules", "./node_modules"))
	assert.True(t, match("/a/**/node_modules", "/a/node_modules"))
	assert.True(t, match("libs/*/*", "libs/common/event-tracking-web-js"))
	assert.False(t, match("libs/*/*", "libs/common/event-tracking-web-js/demo/angular"))
}

func TestMatch(t *testing.T) {
	rootPath := tests.GetRootPath()
	patterns := Patterns2Globs([]string{"**/*.go"})
	ignores := Patterns2Globs([]string{"node_modules", ".git", ".idea", "dist", "**/*_test.go"})
	files := make([]string, 0)
	err := filepath.Walk(rootPath, func(path string, info fs.FileInfo, err error) error {
		rel, err := filepath.Rel(rootPath, path)
		if err != nil {
			return err
		}
		if os.PathSeparator == '\\' {
			rel = strings.ReplaceAll(rel, "\\", "/")
		}
		if MiniMatch(ignores, rel) {
			if info.IsDir() {
				return filepath.SkipDir
			}
			return nil
		}
		if MiniMatch(patterns, rel) {
			files = append(files, rel)
			if info.IsDir() {
				return filepath.SkipDir
			}
			return nil
		}
		return nil
	})
	assert.NoError(t, err)
	_, _ = ffmt.Puts(files)
	assert.NotEmpty(t, files)
	assert.True(t, array.StringEvery(files, func(s string) bool {
		return strings.HasSuffix(s, ".go")
	}))
	assert.False(t, array.StringSome(files, func(s string) bool {
		return strings.HasSuffix(s, "_test.go")
	}))
}
