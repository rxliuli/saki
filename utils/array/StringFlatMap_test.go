package array

import (
	"github.com/stretchr/testify/assert"
	"strings"
	"testing"
)

func TestStringFlatMap(t *testing.T) {
	assert.Equal(t, StringFlatMap([]string{"hello", "world"}, func(s string) []string {
		return strings.Split(s, "")
	}), []string{"h", "e", "l", "l", "o", "w", "o", "r", "l", "d"})
}
