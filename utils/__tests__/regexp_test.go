package __tests__

import (
	"github.com/stretchr/testify/assert"
	"regexp"
	"testing"
)

func TestRegexp(t *testing.T) {
	compile, err := regexp.Compile("\\?raw$")
	assert.NoError(t, err)
	assert.True(t, compile.MatchString("./name.txt?raw"))
}
