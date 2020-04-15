package test

import (
	"fmt"
	"github.com/gitpillow/pfunc"
	"github.com/stretchr/testify/assert"
	"testing"
)

type Person struct {
	Name  string   `json:"name"`
	Age   int      `json:"age"`
	Hobby []string `json:"hobby"`
}

// wrap divide function in python script
func divide(a int, b int) (int, error) {
	i, err := pfunc.Func("dirs/a/b/c/pfunc_test.py", "divide").
		Params(a, b).
		Return(int(0)).
		Do()

	return i.(int), err
}

/// wrap func_return_struct function in python script to assemble a struct
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

/// wrap func_return_struct function in python script to assemble a map with string keys
func GetPersonMap(name string, age int, hobby1 string, hobby2 string) (map[string]interface{}, error) {
	i, err := pfunc.Func("dirs/a/b/c/pfunc_test.py", "func_return_struct").
		Params(name, age, hobby1, hobby2).
		Return(map[string]interface{}{}).
		Do()

	return i.(map[string]interface{}), err
}

func NamesOfThreePeople(ps ...Person) (string, error) {
	i, err := pfunc.Func("dirs/a/b/c/pfunc_test.py", "names_of_three_people").
		ParamDefaults(Person{Name: "SomeOne"}, Person{Name: "OtherOne"}, Person{Name: "AnotherOne"}).
		VarArgs(ps).
		Return("").
		Do()

	return i.(string), err
}

func FirstParamAndOtherParams(first interface{}, name string, age int, hobby []string) (map[string]interface{}, error) {
	r, err := pfunc.Func("dirs/a/b/c/pfunc_test.py", "first_param_and_other_params").
		Params(first).
		KeyWrodParam("name", name).
		KeyWrodParam("age", age).
		KeyWrodParam("hobby", hobby).
		Return(map[string]interface{}{}).
		Do()

	return r.(map[string]interface{}), err
}

func TestWrapFunction(t *testing.T) {
	i, e := divide(6, 3)
	fmt.Println(i)
	fmt.Println(e)
	assert.Equal(t, 2, i)
	assert.Nil(t, e)

	i, e = divide(6, 0)
	fmt.Println(i)
	fmt.Println(e)
	assert.Equal(t, 0, i)
	assert.Contains(t, e.Error(), "by zero")
}

func TestWrapFunctionReturnStruct(t *testing.T) {
	person, err := GetPerson("Tom", 33, "Football", "Shopping")
	fmt.Println(person)
	fmt.Println(err)

	assert.Equal(t,
		Person{
			Name: "Tom",
			Age:  33,
			Hobby: []string{
				"Football",
				"Shopping",
			}},
		person)
	assert.Nil(t, nil)
}

func TestWrapFunctionReturnMap(t *testing.T) {
	personMap, err := GetPersonMap("Tom", 33, "Football", "Shopping")
	fmt.Println(personMap)
	fmt.Println(err)

	assert.Equal(t,
		map[string]interface{}{
			"name": "Tom",
			"age":  float64(33),
			"hobby": []interface{}{
				"Football",
				"Shopping",
			}},
		personMap)
	assert.Nil(t, nil)
}

func TestWrapFunctionWithParamDefaults(t *testing.T) {
	names, err := NamesOfThreePeople()
	fmt.Println(names)
	fmt.Println(err)
	assert.Equal(t, "SomeOne and OtherOne and AnotherOne", names)
	assert.Nil(t, err)

	names, err = NamesOfThreePeople(Person{Name: "Tom"})
	fmt.Println(names)
	fmt.Println(err)
	assert.Equal(t, "Tom and OtherOne and AnotherOne", names)
	assert.Nil(t, err)

	names, err = NamesOfThreePeople(Person{Name: "Tom"}, Person{Name: "Jack"})
	fmt.Println(names)
	fmt.Println(err)
	assert.Equal(t, "Tom and Jack and AnotherOne", names)
	assert.Nil(t, err)
}

func TestWrapFunctionWithKeyWordParams(t *testing.T) {
	r, err := FirstParamAndOtherParams("Lee", "Ming", 33, []string{"Video Game", "Programing"})
	fmt.Println(r)
	fmt.Println(err)
	assert.Equal(t, r, map[string]interface{}{
		"first": "Lee",
		"name":  "Ming",
		"age":   float64(33),
		"hobby": []interface{}{"Video Game", "Programing"},
	})
}
