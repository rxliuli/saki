package fsExtra

import (
	"github.com/stretchr/testify/assert"
	"path/filepath"
	"testing"
)

func TestDirname(t *testing.T) {
	assert.Equal(t, filepath.Base(Dirname()), "fsExtra")
}
