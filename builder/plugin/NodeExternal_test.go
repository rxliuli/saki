package plugin

import (
	"github.com/evanw/esbuild/pkg/api"
	"github.com/rxliuli/saki/utils/fsExtra"
	"github.com/stretchr/testify/assert"
	"path/filepath"
	"testing"
)

func init() {
	beforeEach()
}

func TestNodeExternal(t *testing.T) {
	assert.NoError(t, fsExtra.WriteStringFile(filepath.Join(tempPath, "index.js"), `import * as path from "node:path";

export function hello() {
  console.log(path.resolve());
}
`))
	options := api.BuildOptions{
		EntryPoints: []string{filepath.Join(tempPath, "index.js")},
		Outfile:     filepath.Join(tempPath, "dist/index.js"),
		Bundle:      true,
		Format:      api.FormatESModule,
		Write:       true,
	}
	assert.NotEmpty(t, api.Build(options).Errors)
	options.Plugins = []api.Plugin{
		NodeExternal(),
	}
	assert.Empty(t, api.Build(options).Errors)
}
