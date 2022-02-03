package object

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestAssign(t *testing.T) {
	assert.Equal(t, Assign(map[string]string{
		"name": "琉璃",
		"age":  "18",
	}, map[string]string{
		"name": "灵梦",
	})["name"], "灵梦")
}

//测试会修改原对象
func TestShouldModifyMap(t *testing.T) {
	old := map[string]string{
		"name": "琉璃",
		"age":  "18",
	}
	Assign(old, map[string]string{
		"name": "灵梦",
	})
	assert.Equal(t, old["name"], "灵梦")
}

func TestShouldAssignEmptyMap(t *testing.T) {
	a := map[string]string{"name": "灵梦"}
	b := map[string]string{"age": "17"}
	assert.Equal(t, Assign(map[string]string{}, a, b), map[string]string{"name": "灵梦", "age": "17"})
	assert.Equal(t, a, map[string]string{"name": "灵梦"})
	assert.Equal(t, b, map[string]string{"age": "17"})
}
