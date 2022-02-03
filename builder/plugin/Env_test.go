package plugin

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestDefineImport(t *testing.T) {
	res := loadEnv()
	assert.NotEmpty(t, res["GOENV"])
}
