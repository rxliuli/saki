package __tests__

import (
	"gopkg.in/yaml.v2"
	"strings"
	"testing"
)

func TestYamlParse(t *testing.T) {
	res := struct {
		Packages []string `yaml:"packages"`
	}{}
	_ = yaml.Unmarshal([]byte(`packages:
  - 'libs/*/*'
  - 'apps/*'
  - 'applets/*'
  - 'scripts'
  - 'libs/common/pinefield-cli/templates/*'
  - 'libs/os/sora/templates/sora-app-template-simple'
  - 'libs/os/sora/templates/sora-app-template-vue'
  - '!submodules'
  - '!apps/new-workflow'
  - '!apps/space-console'
  - '!apps/DongAoHuiDataEdit'
`), &res)
	if len(res.Packages) == 0 {
		t.Error("packages is empty")
	}
}

func TestYamlStringify(t *testing.T) {
	marshalStr, err := yaml.Marshal(map[string]interface{}{
		"packages": []string{
			"libs/*/*",
			"apps/*",
			"applets/*",
			"scripts",
			"libs/common/pinefield-cli/templates/*",
			"libs/os/sora/templates/sora-app-template-simple",
			"libs/os/sora/templates/sora-app-template-vue",
			"!submodules",
			"!apps/new-workflow",
			"!apps/space-console",
			"!apps/DongAoHuiDataEdit",
		},
	})
	if err != nil {
		t.Error(err)
	}
	if strings.Index(string(marshalStr), "packages") == -1 {
		t.Error("packages is not in marshalStr")
	}
}
