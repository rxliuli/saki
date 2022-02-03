package plugin

import (
	"github.com/evanw/esbuild/pkg/api"
	"path/filepath"
	"regexp"
)

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
