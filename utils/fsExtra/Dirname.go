package fsExtra

import "path/filepath"

func Dirname() string {
	abs, err := filepath.Abs("./")
	if err != nil {
		panic(err)
	}
	return abs
}
