package fsExtra

import (
	"encoding/json"
	"io/fs"
	"os"
)

func WriteJson(path string, data interface{}) error {
	marshal, err := json.Marshal(data)
	if err != nil {
		return err
	}
	return os.WriteFile(path, marshal, fs.ModePerm)
}
