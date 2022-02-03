package plugin

import (
	"github.com/evanw/esbuild/pkg/api"
	"github.com/rxliuli/saki/utils/jsonext"
	"github.com/rxliuli/saki/utils/object"
	"os"
	"strings"
)

func loadEnv() map[string]string {
	env := map[string]string{}
	for _, v := range os.Environ() {
		kv := strings.SplitN(v, "=", 2)
		if len(kv) == 2 {
			env[kv[0]] = kv[1]
		}
	}
	return env
}

func defineImport(envs map[string]string) map[string]string {
	definitions := map[string]string{}
	for k, v := range envs {
		definitions["import.meta.env."+k] = jsonext.Stringify(v)
	}
	return definitions
}

func defineProcess(envs map[string]string) map[string]string {
	definitions := map[string]string{}
	for k, v := range envs {
		definitions["process.env."+k] = jsonext.Stringify(v)
	}
	return definitions
}

type EnvOptions struct {
	Process bool `json:"process"`
	Import  bool `json:"import"`
}

// Env 环境变量插件
func Env(options EnvOptions) api.Plugin {
	return api.Plugin{
		Name: "esbuild-plugin-env",
		Setup: func(build api.PluginBuild) {
			if build.InitialOptions.Define == nil {
				build.InitialOptions.Define = map[string]string{}
			}
			envs := loadEnv()
			if options.Import {
				build.InitialOptions.Define = object.Assign(build.InitialOptions.Define, defineImport(envs))
			}
			if options.Process {
				build.InitialOptions.Define = object.Assign(build.InitialOptions.Define, defineProcess(envs))
			}
		},
	}
}
