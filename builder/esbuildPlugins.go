package builder

import (
	"github.com/evanw/esbuild/pkg/api"
	"path/filepath"
	"regexp"
	"strings"
)

// NodeExternals 排除和替换 node 内置模块
func NodeExternals() api.Plugin {
	return api.Plugin{
		Name: "esbuild-plugin-node-externals",
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

// AutoExternal 自动排除所有依赖项
func AutoExternal() api.Plugin {
	return api.Plugin{
		Name: "esbuild-plugin-auto-external",
		Setup: func(build api.PluginBuild) {
			compile, err := regexp.Compile("^\\.{1,2}/")
			if err != nil {
				panic("autoExternal regexp compile error")
			}
			build.OnResolve(api.OnResolveOptions{
				Filter: ".*",
			}, func(args api.OnResolveArgs) (api.OnResolveResult, error) {
				if filepath.IsAbs(args.Path) || compile.MatchString(args.Path) {
					return api.OnResolveResult{}, nil
				}
				//fmt.Println("external: ", args.Path)
				return api.OnResolveResult{
					Path:     args.Path,
					External: true,
				}, nil
			})
		},
	}
}
