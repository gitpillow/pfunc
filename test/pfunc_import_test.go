package test

import (
	"fmt"
	"github.com/gitpillow/pfunc"
	"github.com/stretchr/testify/assert"
	"os"
	"path/filepath"
	"testing"
)


func TestImportScriptInSubDir(t *testing.T) {
	result := pfunc.Invoke("dirs/a/b/c/pfunc_test.py", "func2", []interface{}{1, 2})
	assert.Equal(t, true, result.NoError)
	importLine := pfunc.FindLine(result.TempScript, "from", "import")
	fmt.Println(importLine)
	assert.Contains(t, importLine, "from pfunc_test import func2")
}

func TestImportScriptInSubDirWithPythonPath(t *testing.T) {
	appendpp, _ := filepath.Abs("dirs" + string(os.PathSeparator) + "a" + string(os.PathSeparator) + "b")
	env := os.Environ()
	oldpp := os.Getenv(pfunc.PythonPath)
	pfunc.AddEnv(&env, pfunc.PythonPath, appendpp)
	pp, _ := pfunc.GetEnv(&env, pfunc.PythonPath)
	os.Setenv(pfunc.PythonPath, pp)
	newpp := os.Getenv(pfunc.PythonPath)
	fmt.Println(newpp)

	result := pfunc.Invoke("dirs/a/b/c/pfunc_test.py", "func2", []interface{}{1, 2})
	assert.Equal(t, true, result.NoError)
	importLine := pfunc.FindLine(result.TempScript, "from", "import")
	fmt.Println(importLine)
	assert.Contains(t, importLine, "from c.pfunc_test import func2")

	os.Setenv(pfunc.PythonPath, oldpp)
}

func TestImportScriptInSubDirWithPythonPath2(t *testing.T) {
	appendpp, _ := filepath.Abs("dirs" + string(os.PathSeparator) + "a")
	env := os.Environ()
	oldpp := os.Getenv(pfunc.PythonPath)
	pfunc.AddEnv(&env, pfunc.PythonPath, appendpp)
	pp, _ := pfunc.GetEnv(&env, pfunc.PythonPath)

	os.Setenv(pfunc.PythonPath, pp)
	newpp := os.Getenv(pfunc.PythonPath)
	fmt.Println(newpp)

	result := pfunc.Invoke("dirs/a/b/c/pfunc_test.py", "func2", []interface{}{1, 2})
	assert.Equal(t, true, result.NoError)
	importLine := pfunc.FindLine(result.TempScript, "from", "import")
	fmt.Println(importLine)
	assert.Contains(t, importLine, "from b.c.pfunc_test import func2")

	os.Setenv(pfunc.PythonPath, oldpp)
}
