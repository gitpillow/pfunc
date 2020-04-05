package test

import (
	"fmt"
	"github.com/gitpillow/pfunc"
	"testing"
)

func TestInvoke(t *testing.T) {
	result := pfunc.Invoke("./pfunc_test.py", "func1", nil)
	fmt.Println(result)

	result = pfunc.Invoke("pfunc_test.py", "func2", []interface{}{1, 2})
	fmt.Println(result)

	result = pfunc.Invoke("pfunc_test.py", "func3", []interface{}{
		map[string]interface{}{"name": "Jack"},
		map[string]interface{}{"name": "Alice"},
		map[string]interface{}{"name": "Tom"},
	})
	fmt.Println(result)

	result = pfunc.Invoke("pfunc_test.py", "func4", []interface{}{2, 0})
	fmt.Println(result)
}
