package plugin

import (
	"github.com/evanw/esbuild/pkg/api"
	"strings"
)

// NodeExternal 排除和替换 node 内置模块
func NodeExternal() api.Plugin {
	return api.Plugin{
		Name: "esbuild-plugin-node-external",
		Setup: func(build api.PluginBuild) {
			build.OnResolve(api.OnResolveOptions{
				Filter: `^node:`,
			}, func(args api.OnResolveArgs) (api.OnResolveResult, error) {
				return api.OnResolveResult{
					Path:     strings.TrimPrefix(args.Path, "node:"),
					External: true,
				}, nil
			})
		},
	}

}
