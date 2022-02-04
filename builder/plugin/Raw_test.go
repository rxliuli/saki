package plugin

import (
	"github.com/evanw/esbuild/pkg/api"
	"github.com/rxliuli/saki/utils/fsExtra"
	"github.com/stretchr/testify/assert"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

var tempPath = filepath.Join(fsExtra.Dirname(), ".temp")

func beforeEach() {
	_ = os.RemoveAll(tempPath)
	_ = os.MkdirAll(tempPath, fs.ModeDir)
}

func init() {
	beforeEach()
}

func TestRaw(t *testing.T) {
	assert.NoError(t, fsExtra.WriteStringFile(filepath.Join(tempPath, "name.txt"), "liuli"))
	assert.NoError(t, fsExtra.WriteStringFile(filepath.Join(tempPath, "index.js"), `import name from "./name.txt?raw";

export function hello() {
  console.log("hello " + name);
}`))
	result := api.Build(api.BuildOptions{
		EntryPoints: []string{filepath.Join(tempPath, "index.js")},
		Outfile:     filepath.Join(tempPath, "dist/index.js"),
		Bundle:      true,
		Plugins: []api.Plugin{
			Raw(),
		},
	})
	assert.Empty(t, result.Errors)
	assert.True(t, strings.Contains(string(result.OutputFiles[0].Contents), "liuli"))
}
