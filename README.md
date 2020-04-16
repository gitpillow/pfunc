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

```python
def add(a, b):
    return a + b
```

```go
func divide(a int, b int) (int, error) {
	i, err := pfunc.Func("dirs/a/b/c/pfunc_test.py", "divide").
		Params(a, b).
		Return(int(0)).
		Do()

	return i.(int), err
}

i, e := divide(6, 3)
fmt.Println(i)
fmt.Println(e)
```

output:
```shell script
2
<nil>
```

## configure

