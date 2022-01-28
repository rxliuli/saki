package build

import (
	"github.com/evanw/esbuild/pkg/api"
	"path/filepath"
	"regexp"
	"strings"
)

func NativeNodeModules() api.Plugin {
	return api.Plugin{Name: "native-node-modules", Setup: func(build api.PluginBuild) {
		namespace := "node-file"
		build.OnResolve(api.OnResolveOptions{
			Filter:    "\\.node$",
			Namespace: "file",
		}, func(args api.OnResolveArgs) (api.OnResolveResult, error) {
			return api.OnResolveResult{
				Path:      args.Path,
				Namespace: namespace,
			}, nil
		})
		build.OnLoad(api.OnLoadOptions{
			Filter:    ".*",
			Namespace: namespace,
		}, func(args api.OnLoadArgs) (api.OnLoadResult, error) {
			contents := `
        import path from ${JSON.stringify(args.path)}
        try { module.exports = require(path) }
        catch {}
      `
			return api.OnLoadResult{
				Contents: &contents,
			}, nil
		})
		build.OnResolve(api.OnResolveOptions{
			Filter:    "\\.node$",
			Namespace: namespace,
		}, func(args api.OnResolveArgs) (api.OnResolveResult, error) {
			return api.OnResolveResult{
				Path:      args.Path,
				Namespace: "file",
			}, nil
		})
		options := build.InitialOptions
		if options.Loader == nil {
			options.Loader = map[string]api.Loader{}
		}
		options.Loader[".node"] = api.LoaderFile
	}}
}

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
			compile, err := regexp.Compile("^\\.{1,2}\\/")
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
