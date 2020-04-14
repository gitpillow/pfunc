package test

import (
	"fmt"
	"github.com/gitpillow/pfunc"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestInvokeWithOutReturn(t *testing.T) {
	result := pfunc.Invoke("dirs/a/b/c/pfunc_test.py", "do_print", nil)

	fmt.Println(result.Inspect())
	assert.Equal(t, true, result.NoError)
	assert.NotEmpty(t, result.TempScript)
	assert.Empty(t, result.Exception)
	assert.Equal(t, "null", result.JsonRepresentation)
}

func TestInvokeWithIntReturn(t *testing.T) {
	result := pfunc.Invoke("dirs/a/b/c/pfunc_test.py", "add", []interface{}{1, 2})

	fmt.Println(result.Inspect())
	assert.Equal(t, true, result.NoError)
	assert.Empty(t, result.Exception)
	assert.Equal(t, 3, result.MustInt())
}

func TestInvokeWithStringReturn(t *testing.T) {
	result := pfunc.Invoke("dirs/a/b/c/pfunc_test.py", "names_of_three_people", []interface{}{
		map[string]interface{}{"name": "Jack"},
		map[string]interface{}{"name": "Alice"},
		map[string]interface{}{"name": "Tom"},
	})

	fmt.Println(result.Inspect())
	assert.Equal(t, true, result.NoError)
	assert.Empty(t, result.Exception)
	assert.Equal(t, "Jack and Alice and Tom", result.MustString())
}

func TestInvokeWithException(t *testing.T) {
	result := pfunc.Invoke("dirs/a/b/c/pfunc_test.py", "divide", []interface{}{2, 0})

	fmt.Println(result.Inspect())
	assert.Equal(t, false, result.NoError)
	assert.NotEmpty(t, result.Exception)
	assert.Empty(t, result.JsonRepresentation)
	assert.Equal(t, "", result.JsonRepresentation)
}

func TestInvokeWithFloatReturn(t *testing.T) {
	result := pfunc.Invoke("dirs/a/b/c/pfunc_test.py", "float_divide", []interface{}{3.0, 2.0})

	fmt.Println(result.Inspect())
	assert.Equal(t, true, result.NoError)
	assert.Equal(t, float32(1.5), result.MustFloat())
}
