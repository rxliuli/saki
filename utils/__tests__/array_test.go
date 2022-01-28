package __tests__

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestRemove(t *testing.T) {
	ints := []int{1, 2, 3, 4, 5}
	i := 1
	ints = append(ints[:i], ints[i+1:]...)
	assert.Equal(t, ints, []int{1, 3, 4, 5})
}
