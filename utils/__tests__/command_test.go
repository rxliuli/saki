package __tests__

import (
	"fmt"
	"github.com/rxliuli/saki/tests"
	"github.com/stretchr/testify/assert"
	"os/exec"
	"testing"
)

func init() {
	rootDir = tests.GetRootPath()
}

func TestExecCmd(t *testing.T) {
	command := exec.Command("go", "run", "main.go", "-h")
	command.Dir = rootDir
	output, err := command.Output()
	assert.NoError(t, err)
	fmt.Printf("%s\n", output)
}
