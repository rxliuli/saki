package jsonext

import (
	"encoding/json"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestStringify(t *testing.T) {
	v := map[string]interface{}{
		"a": 1,
		"b": "2",
		"c": []interface{}{1, 2, 3},
	}
	res := map[string]interface{}{}
	assert.NoError(t, json.Unmarshal([]byte(Stringify(v)), &res))
	assert.Equal(t, Stringify(v), Stringify(res))
}
