package test

import (
	"fmt"
	"github.com/gitpillow/pfunc"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestInvokeWithOutReturn(t *testing.T) {
	result := pfunc.Invoke("./pfunc_test.py", "func1", nil)
	fmt.Println(result.Inspect())
	assert.Equal(t, true, result.NoError)
	assert.NotEmpty(t, result.TempScript)
	assert.Empty(t, result.Exception)
	assert.Equal(t, "null", result.JsonRepresentation)
}

func TestInvokeWithIntReturn(t *testing.T) {
	result := pfunc.Invoke("pfunc_test.py", "func2", []interface{}{1, 2})
	fmt.Println(result.Inspect())
	assert.Equal(t, true, result.NoError)
	assert.NotEmpty(t, result.TempScript)
	assert.Empty(t, result.Exception)
	assert.Equal(t, 3, result.MustInt())
}

func TestInvokeWithStringReturn(t *testing.T) {
	result := pfunc.Invoke("pfunc_test.py", "func3", []interface{}{
		map[string]interface{}{"name": "Jack"},
		map[string]interface{}{"name": "Alice"},
		map[string]interface{}{"name": "Tom"},
	})
	fmt.Println(result.Inspect())
	assert.Equal(t, true, result.NoError)
	assert.NotEmpty(t, result.TempScript)
	assert.Empty(t, result.Exception)
	assert.Equal(t, "Jack and Alice and Tom", result.MustString())
}

func TestInvokeWithException(t *testing.T) {
	result := pfunc.Invoke("pfunc_test.py", "func4", []interface{}{2, 0})
	fmt.Println(result.Inspect())
	assert.Equal(t, false, result.NoError)
	assert.NotEmpty(t, result.TempScript)
	assert.NotEmpty(t, result.Exception)
	assert.Empty(t, result.JsonRepresentation)
	assert.Equal(t, "", result.JsonRepresentation)
}
