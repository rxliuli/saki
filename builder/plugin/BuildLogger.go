package plugin

import (
	"fmt"
	"github.com/evanw/esbuild/pkg/api"
	"path/filepath"
	"strings"
	"time"
)

func BuildLogger(cwd string) api.Plugin {
	return api.Plugin{
		Name: "esbuild-plugin-build-logger",
		Setup: func(build api.PluginBuild) {
			start := time.Now()
			build.OnStart(func() (api.OnStartResult, error) {
				start = time.Now()
				return api.OnStartResult{}, nil
			})
			build.OnEnd(func(result *api.BuildResult) {
				rel, _ := filepath.Rel(cwd, build.InitialOptions.Outfile)
				if len(result.Errors) != 0 {
					message := result.Errors[0]
					location := message.Location
					fmt.Printf("构建失败 %s: %s %s:%d:%d\n", rel, message.Text, location.File, location.Line, location.Column)
					return
				}
				fmt.Printf("构建成功 %s: %v\n", strings.ReplaceAll(rel, "\\", "/"), time.Since(start))
			})
		},
	}
}
