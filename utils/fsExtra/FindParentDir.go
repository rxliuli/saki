package fsExtra

import "path/filepath"

// FindParentDir 向上查找满足条件的目录
func FindParentDir(dir string, fn func(dir string) bool) string {
	var prev = ""
	var current = dir
	for current != prev {
		if fn(current) {
			return current
		}
		prev = current
		current = filepath.Dir(current)
	}
	return ""
}
