package fsExtra

import (
	"encoding/json"
	"os"
)

func ReadJson(filePath string, res interface{}) error {
	file, err := os.ReadFile(filePath)
	if err != nil {
		return err
	}
	err = json.Unmarshal(file, &res)
	if err != nil {
		return err
	}
	return nil
}
