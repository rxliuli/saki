package fsExtra

import (
	"path/filepath"
	"testing"
)

func TestDirname(t *testing.T) {
	if filepath.Base(Dirname()) != "fsExtra" {
		t.Errorf("Dirname() failed")
	}
}
