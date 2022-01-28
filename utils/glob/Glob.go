package glob

import (
	"os"
	"path"
	"path/filepath"
	"strings"
)

type Options struct {
	Cwd    string
	Ignore []string
}

func match(patterns []string, name string) bool {
	for _, pattern := range patterns {
		match, err := path.Match(pattern, name)
		if err != nil {
			return false
		}
		if match {
			return true
		}
	}
	return false
}

// Glob 按照指定模式匹配文件或目录
func Glob(patterns []string, options Options) ([]string, error) {
	files := make([]string, 0)
	err := filepath.Walk(options.Cwd, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		realPath, err := filepath.Rel(options.Cwd, path)
		if err != nil {
			return err
		}
		if realPath == "." {
			return nil
		}
		realPath = strings.ReplaceAll(realPath, "\\", "/")
		if match(options.Ignore, realPath) {
			return nil
		}
		//println(realPath)
		if match(patterns, realPath) {
			files = append(files, realPath)
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	return files, nil
}
