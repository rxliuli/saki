package fsExtra

import (
	"os"
)

func ReadStringFile(filePath string) (error, string) {
	file, err := os.ReadFile(filePath)
	if err != nil {
		return err, ""
	}
	return nil, string(file)
}
