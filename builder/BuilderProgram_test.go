package builder

import (
	"github.com/evanw/esbuild/pkg/api"
	"github.com/rxliuli/saki/utils/fsExtra"
	"github.com/stretchr/testify/assert"
	"io/fs"
	"os"
	"path/filepath"
	"testing"
)

var builder = Program{
	Cwd:   filepath.Join("C:/Users/rxliuli/Code/Web/liuli-tools/apps/liuli-cli"),
	Watch: false,
}

var tempPath = filepath.Join(fsExtra.Dirname(), ".temp")

func beforeEach() {
	builder.Cwd = tempPath
	_ = os.RemoveAll(tempPath)
	_ = os.MkdirAll(tempPath, fs.ModeDir)
}

func TestBuilderProgram_BuildLib(t *testing.T) {
	assert.NoError(t, builder.BuildToTargets([]Target{TargetEsm, TargetCjs}))
}

func TestBuilderProgram_BuildCli(t *testing.T) {
	assert.NoError(t, builder.BuildToTargets([]Target{TargetCli, TargetEsm, TargetCjs}))
}

func TestGetPlatformOfNode(t *testing.T) {
	beforeEach()
	_ = fsExtra.WriteJson(filepath.Join(tempPath, "package.json"), map[string]interface{}{
		"name":    "test",
		"version": "1.0.0",
		"devDependencies": map[string]interface{}{
			"@types/node": "*",
		},
	})
	assert.Equal(t, builder.getPlatform(), api.PlatformNode)
}

func TestGetPlatformOfBrowser(t *testing.T) {
	beforeEach()
	_ = fsExtra.WriteStringFile(filepath.Join(tempPath, "tsconfig.json"), `{
  "compilerOptions": {
    //依赖的内置模块
    "lib": [
      "dom",
      "esnext"
    ]
  }
}`)
	assert.Equal(t, builder.getPlatform(), api.PlatformBrowser)
}

func TestGetPlatformOfNeutral(t *testing.T) {
	beforeEach()
	assert.Equal(t, builder.getPlatform(), api.PlatformNeutral)
}

func TestProgram_BuildLibWatch(t *testing.T) {
	beforeEach()
	assert.NoError(t, os.MkdirAll(filepath.Join(tempPath, "src"), fs.ModeDir))
	assert.NoError(t, fsExtra.WriteStringFile(filepath.Join(tempPath, "package.json"), `{}`))
	assert.NoError(t, fsExtra.WriteStringFile(filepath.Join(tempPath, "src/index.ts"), `export function hello(name: string) {
  return 'hello ' + name;
}`))
	builder.Watch = true
	builder.Cwd = tempPath
	assert.NoError(t, builder.BuildToTargets([]Target{TargetEsm, TargetCjs}))
	<-make(chan bool)
}
