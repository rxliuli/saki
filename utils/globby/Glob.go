package globby

import (
	"github.com/rxliuli/saki/utils/array"
	"path/filepath"
)

type Options struct {
	Cwd    string
	Ignore []string
}

// Glob 按照指定模式匹配文件或目录
func Glob(patterns []string, options Options) []string {
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
