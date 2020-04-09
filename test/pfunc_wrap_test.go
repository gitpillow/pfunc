package test

import (
	"fmt"
	"github.com/gitpillow/pfunc"
	"testing"
)

type Person struct {
	Name  string   `json:"name"`
	Age   int      `json:"age"`
	Hobby []string `json:"hobby"`
}

func divide(a int, b int) (int, error) {
	i, err := pfunc.Func("dirs/a/b/c/pfunc_test.py", "func4").
		Params(a, b).
		Return(int(0)).
		Do()

	return i.(int), err
}

func GetPerson(name string, age int, hobby1 string, hobby2 string) (Person, error) {
	i, err := pfunc.Func("dirs/a/b/c/pfunc_test.py", "func_return_struct").
		Params(name, age, hobby1, hobby2).
		Return(Person{
			Name:  "Nobody",
			Age:   0,
			Hobby: []string{},
		}).
		Do()

	return i.(Person), err
}

func GetPersonMap(name string, age int, hobby1 string, hobby2 string) (map[string]interface{}, error) {
	i, err := pfunc.Func("dirs/a/b/c/pfunc_test.py", "func_return_struct").
		Params(name, age, hobby1, hobby2).
		Return(map[string]interface{}{}).
		Do()

	return i.(map[string]interface{}), err
}

func TestWrapFunction(t *testing.T) {
	i, e := divide(6, 3)
	fmt.Println(i)
	fmt.Println(e)

	i, e = divide(6, 0)
	fmt.Println(i)
	fmt.Println(e)
}

func TestWrapFunctionReturnStruct(t *testing.T) {
	person, err := GetPerson("Tom", 33, "Football", "Shopping")
	fmt.Println(person)
	fmt.Println(err)
}

func TestWrapFunctionReturnMap(t *testing.T) {
	personMap, err := GetPersonMap("Tom", 33, "Football", "Shopping")
	fmt.Println(personMap)
	fmt.Println(err)
}
