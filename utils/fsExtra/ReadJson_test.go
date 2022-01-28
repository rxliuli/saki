package fsExtra

import (
	"github.com/swaggest/assertjson/json5"
	"io/fs"
	"os"
	"path/filepath"
	"testing"
)

var tempPath = filepath.Join(Dirname(), ".temp")

func init() {
	_ = os.RemoveAll(tempPath)
	_ = os.MkdirAll(tempPath, fs.ModeDir)
}

func TestReadJson(t *testing.T) {
	jsonPath := filepath.Join(tempPath, "test.json")
	_ = WriteStringFile(jsonPath, `{"name": "test"}`)
	var res struct {
		Name string `json:"name"`
	}
	_ = ReadJson(jsonPath, &res)
	if res.Name != "test" {
		t.Error("read json failed")
	}
}

func TestReadJson5(t *testing.T) {
	jsonPath := filepath.Join(tempPath, "test.json")
	_ = WriteStringFile(jsonPath, `{
  //这是注释
  name: "test",
}`)
	s, err := os.ReadFile(jsonPath)
	if err != nil {
		t.Error("read file failed")
	}
	var res struct {
		Name string `json:"name"`
	}
	_ = json5.Unmarshal(s, &res)
	if res.Name != "test" {
		t.Error("read json failed")
	}
}
