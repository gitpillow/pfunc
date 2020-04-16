# pfunc
pfunc is a tiny go lib to call python funcions in python scrips. Relative to raw command execution, it behaves better in arguments passing, 
return value acquisition, and exception handling. 

## usage

### simple demo
It is easy to invoke a function defined in a python script.

pfunc_test.py
```python
def add(a, b):
    return a + b
```

Go code
```go
result := pfunc.Invoke("dirs/a/b/c/pfunc_test.py", "add", []interface{}{1, 2})
fmt.Println(result.Inspect())
```

As Inspect() info showed before, Inovke() return An PResult struct instance which has some field to tell the function 
execution detail and result. 

```go
result.NoError                  // tell the function execution is success of fail with exception
result.JsonRepresentation       // serialize function return value to json string 
result.Exception                // python error print string
result.TempScript               // full text of temp script which has been execution
result.Output                   // temp script execution output
result.PythonPath               // PYTHONPATH env value of current execution
```

output
```shell script
python function result: 
    status:
        success
    return value:
        3
    exception:

    temp script:
        
        from pfunc_test import add
        import traceback
        import json
        try:
            null = None
            pfunc_inject_0 = 1
            pfunc_inject_1 = 2
            result = add(pfunc_inject_0, pfunc_inject_1)
            print 'pfunc_return_start_{}pfunc_return_end_'.format(json.dumps(result))
        except Exception, e:
            msg = traceback.format_exc()
            print "pfunc_exception_start_",
            print msg,
            print "pfunc_exception_end_",
    python path:
        D:\projects\pfunc\test\dirs\a\b\c
```

### wrap python function as go function

It is simple to wrap an python function, too

#### function with position params

pfunc_test.py
```python
def add(a, b):
    return a + b
```

Go code
```go
func divide(a int, b int) (int, error) {
	i, err := pfunc.Func("dirs/a/b/c/pfunc_test.py", "divide").             // set script path, set function name
		Params(a, b).                                                       // set parameters type and value
		Return(int(0)).                                                     // set return type and default return value
		Do()                                                                // invoke

	return i.(int), err                                                     // type cast and return
}

i, e := divide(6, 3)
fmt.Println(i)
fmt.Println(e)
if e != nil {
    // error handle
}
```

output:
```shell script
2
<nil>
```

#### function return struct or map[string]interface{}

pfunc_test.py
```python
def func_return_struct(name, age, hobby1, hobby2):
    return {
        "name": name,
        "age": age,
        "hobby": [
            hobby1,
            hobby2
        ]
    }
```

Go code
```go
func GetPerson(name string, age int, hobby1 string, hobby2 string) (Person, error) {
	i, err := pfunc.Func("dirs/a/b/c/pfunc_test.py", "func_return_struct").                 // set script path , set function name
		Params(name, age, hobby1, hobby2).                                                  // set parameter types and values
		Return(Person{                                                                      // set return type and default return value
			Name:  "Nobody",
			Age:   0,
			Hobby: []string{},
		}).
		Do()                                                                                // invoke

	return i.(Person), err                                                                  // type cast and return
}

/// wrap func_return_struct function in python script to assemble a map with string keys
func GetPersonMap(name string, age int, hobby1 string, hobby2 string) (map[string]interface{}, error) {
	i, err := pfunc.Func("dirs/a/b/c/pfunc_test.py", "func_return_struct").                 // set script path, set function name
		Params(name, age, hobby1, hobby2).                                                  // set parameter types and values
		Return(map[string]interface{}{}).                                                   // set return type and default return value
		Do()                                                                                // invoke

	return i.(map[string]interface{}), err                                                  // type cast and return
}

person, err := GetPerson("Tom", 33, "Football", "Shopping")
fmt.Println(person)
fmt.Println(err)

personMap, err := GetPersonMap("Tom", 33, "Football", "Shopping")
fmt.Println(personMap)
fmt.Println(err)

```

output:
```shell script
{Tom 33 [Football Shopping]}
<nil>

map[age:33 hobby:[Football Shopping] name:Tom]
<nil>
```


## configure

