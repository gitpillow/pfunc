package test

import (
	"fmt"
	"github.com/gitpillow/pfunc"
	"testing"
)

func divide(a int, b int) (int, error) {
	i, err := pfunc.Func("dirs/a/b/c/pfunc_test.py", "func4").
		Params(a, b).
		Return(int(0)).
		Do()

	return i.(int), err
}

func TestWrapFunction(t *testing.T) {
	i, e := divide(6, 3)
	fmt.Println(i)
	fmt.Println(e)

	i, e = divide(6, 0)
	fmt.Println(i)
	fmt.Println(e)
}
