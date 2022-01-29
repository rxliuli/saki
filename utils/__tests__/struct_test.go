package __tests__

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

type Person struct {
	Name string
}

func TestStruct(t *testing.T) {
	persons := []Person{{Name: "琉璃"}}
	person := persons[0]
	person.Name = "灵梦"
	assert.NotEqual(t, persons[0].Name, person.Name)
}
