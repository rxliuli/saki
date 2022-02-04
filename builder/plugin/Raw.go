package plugin

import (
	"github.com/evanw/esbuild/pkg/api"
	"github.com/rxliuli/saki/utils/fsExtra"
	"path/filepath"
	"strings"
)

func Raw() api.Plugin {
	return api.Plugin{
		Name: "esbuild-plugin-raw",
		Setup: func(build api.PluginBuild) {
			build.OnResolve(api.OnResolveOptions{
				Filter: "\\?raw$",
			}, func(args api.OnResolveArgs) (api.OnResolveResult, error) {
				var path string
				if filepath.IsAbs(args.Path) {
					path = args.Path
				} else {
					path = filepath.Join(args.ResolveDir, args.Path)
				}
				return api.OnResolveResult{
					Path:      path,
					Namespace: "raw-loader",
				}, nil
			})
			build.OnLoad(api.OnLoadOptions{
				Filter:    "\\?raw$",
				Namespace: "raw-loader",
			}, func(args api.OnLoadArgs) (api.OnLoadResult, error) {
				err, s := fsExtra.ReadStringFile(strings.TrimRight(args.Path, "?raw"))
				if err != nil {
					return api.OnLoadResult{}, err
				}
				return api.OnLoadResult{
					Contents: &s,
					Loader:   api.LoaderText,
				}, nil
			})
		},
	}
}
