package globby

import (
	"github.com/rxliuli/saki/utils/array"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
)

type Options struct {
	Cwd    string
	Ignore []string
	First  bool
}

// Glob 按照指定模式匹配文件或目录
func simpleGlob(patterns []string, options Options) []string {
	files := make([]string, 0)
	for _, pattern := range patterns {
		matches, err := filepath.Glob(filepath.Join(options.Cwd, pattern))
		if err != nil {
			continue
		}
		files = append(files, matches...)
	}
	return array.StringMap(files, func(file string) string {
		rel, err := filepath.Rel(options.Cwd, file)
		if err != nil {
			panic("filepath.Rel error")
		}
		return rel
	})
}

func Glob(ps []string, options Options) []string {
	rootPath := options.Cwd
	patterns := Patterns2Globs(ps)
	ignores := Patterns2Globs(options.Ignore)
	files := make([]string, 0)
	err := filepath.Walk(rootPath, func(path string, info fs.FileInfo, err error) error {
		rel, err := filepath.Rel(rootPath, path)
		if rel == "." {
			return nil
		}
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
			//if info.IsDir() {
			//	return filepath.SkipDir
			//}
			return nil
		}
		return nil
	})
	if err != nil {
		panic(err)
	}
	return files
}
