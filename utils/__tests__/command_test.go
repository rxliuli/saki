package __tests__

import (
	"fmt"
	"github.com/rxliuli/saki/utils/fsExtra"
	"github.com/stretchr/testify/assert"
	"os/exec"
	"path/filepath"
	"testing"
)

func init() {
	rootDir = fsExtra.FindParentDir(fsExtra.Dirname(), func(dir string) bool {
		return fsExtra.PathExists(filepath.Join(dir, "go.mod"))
	})
}

func TestExecCmd(t *testing.T) {
	command := exec.Command("go", "run", "main.go", "-h")
	command.Dir = rootDir
	output, err := command.Output()
	assert.NoError(t, err)
	fmt.Printf("%s\n", output)
}
