package __tests__

import (
	"github.com/stretchr/testify/assert"
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
	assert.NotEmpty(t, res.Packages)
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
	assert.NoError(t, err)
	assert.Equal(t, strings.Index(string(marshalStr), "packages"), 0)
}
