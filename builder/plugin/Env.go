package plugin

import (
	"github.com/evanw/esbuild/pkg/api"
	"github.com/rxliuli/saki/utils/jsonext"
	"github.com/rxliuli/saki/utils/object"
	"os"
	"regexp"
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
	compile, err := regexp.Compile("^[_$a-zA-Z][$_a-zA-Z0-9]*$")
	if err != nil {
		panic(err)
	}
	for k, v := range envs {
		if compile.Match([]byte(k)) {
			definitions["import.meta.env."+k] = jsonext.Stringify(v)
		}
	}
	return definitions
}

//func defineProcess(envs map[string]string) map[string]string {
//	definitions := map[string]string{}
//	for k, v := range envs {
//		definitions["process.env."+k] = jsonext.Stringify(v)
//	}
//	return definitions
//}

// Env 环境变量插件
func Env() api.Plugin {
	return api.Plugin{
		Name: "esbuild-plugin-env",
		Setup: func(build api.PluginBuild) {
			if build.InitialOptions.Define == nil {
				build.InitialOptions.Define = map[string]string{}
			}
			build.InitialOptions.Define = object.Assign(build.InitialOptions.Define, defineImport(loadEnv()))
		},
	}
}
