package __tests__

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
	"time"
)

func TestRemoveAll(t *testing.T) {
	start := time.Now()
	err := os.RemoveAll("C:/Users/rxliuli/Code/company/matrix/node_modules")
	fmt.Println("removeAll:", time.Now().Sub(start))
	assert.NoError(t, err)
}
