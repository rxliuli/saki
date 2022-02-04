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

func TestAutoExternal(t *testing.T) {
	assert.NoError(t, fsExtra.WriteStringFile(filepath.Join(tempPath, "index.js"), `import { uniqBy } from "lodash";

export function hello() {
  console.log(uniqBy([1, 2, 3, 4, 5, 6, 7, 8, 9, 10]));
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
		AutoExternal(),
	}
	assert.Empty(t, api.Build(options).Errors)
}
