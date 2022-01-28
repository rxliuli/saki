package fsExtra

import (
	"github.com/stretchr/testify/assert"
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
	assert.Equal(t, res.Name, "test")
}

func TestReadJson5(t *testing.T) {
	jsonPath := filepath.Join(tempPath, "test.json")
	_ = WriteStringFile(jsonPath, `{
  //这是注释
  name: "test",
}`)
	s, err := os.ReadFile(jsonPath)
	assert.NoError(t, err)
	var res struct {
		Name string `json:"name"`
	}
	_ = json5.Unmarshal(s, &res)
	assert.Equal(t, res.Name, "test")
}
