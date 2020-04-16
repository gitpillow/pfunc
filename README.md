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

## configure

