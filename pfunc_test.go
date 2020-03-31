package pfunc

import (
	"testing"
)

func TestInvoke(t *testing.T) {
	result, script, err := Invoke("pfunc_test.py", "func1", nil)
	printInfo(t, result, script, err)

	result, script, err = Invoke("pfunc_test.py", "func2", []interface{}{1, 2})
	printInfo(t, result, script, err)

	result, script, err = Invoke("pfunc_test.py", "func3", []interface{}{
		map[string]interface{}{"name": "Jack"},
		map[string]interface{}{"name": "Alice"},
		map[string]interface{}{"name": "Tom"},
	})
	printInfo(t, result, script, err)

	result, script, err = Invoke("pfunc_test.py", "func4", []interface{}{2, 0})
	printInfo(t, result, script, err)
}

func printInfo(t *testing.T, result interface{}, script string, err error) {
	t.Log("---- script is:")
	t.Log(script)
	t.Log("---- result is:")
	t.Log(result)
	t.Log("---- exception is:")
	t.Log(err)
}
