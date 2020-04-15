package pfunc

import (
	"bufio"
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"reflect"
	"strconv"
	"strings"
)

const InjectVarNamePrefix string = "pfunc_inject_"
const ReturnValueStart string = "pfunc_return_start_"
const ReturnValueEnd string = "pfunc_return_end_"
const ExceptionStart string = "pfunc_exception_start_"
const ExceptionEnd string = "pfunc_exception_end_"
const PythonExecutable string = "python"
const PythonPath string = "PYTHONPATH"

const PResultToString = `
python function result: 
    status:
        %v
    return value:
%v
    exception:
%v
    temp script:
%v
    python path:
        %v
`

const Python2ScriptTemplate string = `
from %s import %s
import traceback
import json
try:
    null = None
%s
    result = %s
    print '%s{}%s'.format(json.dumps(result))
except Exception, e:
    msg = traceback.format_exc()
    print "%s",
    print msg,
    print "%s",
`

// invoke result struct
type PResult struct {
	NoError            bool
	JsonRepresentation string
	Exception          error
	TempScript         string
	Output             string
	PythonPath         string
}

type WrapInfo struct {
	scriptPath         string
	funcName           string
	returnType         reflect.Type
	returnValue        interface{}
	paramTypes         []reflect.Type
	paramValues        []interface{}
	paramDefaultValues []interface{}
	Keywords           map[string]interface{}
	wrapError          []error
}

func (pr PResult) Inspect() string {
	return fmt.Sprintf(PResultToString,
		Select(pr.NoError, "success", "fail").(string),
		TabString(pr.JsonRepresentation, 8),
		TabString(pr.Exception.Error(), 8),
		TabString(pr.TempScript, 8),
		pr.PythonPath)
}

func (pr PResult) Int() (int, error) {
	return strconv.Atoi(pr.JsonRepresentation)
}

func (pr PResult) MustInt() int {
	i, _ := pr.Int()
	return i
}

func (pr PResult) String() (string, error) {
	var str string
	err := json.Unmarshal([]byte(pr.JsonRepresentation), &str)
	return str, err
}

func (pr PResult) MustString() string {
	s, _ := pr.String()
	return s
}

func (pr PResult) Float() (float32, error) {
	var f float32
	err := json.Unmarshal([]byte(pr.JsonRepresentation), &f)
	return f, err
}

func (pr PResult) MustFloat() float32 {
	f, _ := pr.Float()
	return f
}

func Func(scriptPath string, funcName string) *WrapInfo {
	wrapInfo := &WrapInfo{}
	wrapInfo.scriptPath = scriptPath
	wrapInfo.funcName = funcName
	return wrapInfo
}

func (w *WrapInfo) Return(i interface{}) *WrapInfo {
	t := reflect.TypeOf(i)
	w.returnType = t
	w.returnValue = i
	return w
}

func (w *WrapInfo) Params(interfaces ...interface{}) *WrapInfo {
	w.paramTypes = []reflect.Type{}
	w.paramValues = interfaces
	for _, i := range interfaces {
		w.paramTypes = append(w.paramTypes, reflect.TypeOf(i))
	}
	return w
}

func (w *WrapInfo) KeyWrodParam(keyword string, value interface{}) *WrapInfo {
	if w.Keywords == nil {
		w.Keywords = map[string]interface{}{}
	}
	w.Keywords[keyword] = value
	return w
}

func (w *WrapInfo) ParamDefaults(interfaces ...interface{}) *WrapInfo {
	w.paramDefaultValues = interfaces
	return w
}

func (w *WrapInfo) VarArgs(varargs interface{}) *WrapInfo {
	t := reflect.TypeOf(varargs)
	k := t.Kind()
	if k != reflect.Slice {
		w.wrapError = append(w.wrapError, fmt.Errorf("varargs is not of slice type"))
	} else {
		var params []interface{}
		v := reflect.ValueOf(varargs)
		for i := 0; i < v.Len(); i++ {
			e := v.Index(i)
			params = append(params, e.Interface())
		}
		w.Params(params...)
	}
	return w
}

func (w *WrapInfo) Do(interfaces ...interface{}) (interface{}, error) {

	if w.returnType == nil {
		return nil, fmt.Errorf("return type is not set")
	}

	if len(w.wrapError) > 0 {
		return nil, w.wrapError[0]
	}

	if len(w.paramDefaultValues) > 0 {
		for i, d := range w.paramDefaultValues {
			if i >= len(w.paramValues) {
				w.paramValues = append(w.paramValues, d)
			}
		}
	}

	r := doInvoke(w.scriptPath, w.funcName, w.paramValues, w.Keywords)
	if r.NoError {
		i := reflect.New(w.returnType).Interface()
		err := json.Unmarshal([]byte(r.JsonRepresentation), i)
		if err != nil {
			return w.returnValue, err
		}
		return reflect.ValueOf(i).Elem().Interface(), nil
	} else {
		e := reflect.New(w.returnType).Elem()
		e.Set(reflect.ValueOf(w.returnValue))
		return e.Interface(), r.Exception
	}
}

func Call(scriptPath string, funcName string, params ...interface{}) PResult {
	return Invoke(scriptPath, funcName, params)
}

func Invoke(scriptPath string, funcName string, params []interface{}) PResult {
	return doInvoke(scriptPath, funcName, params, nil)
}

func doInvoke(scriptPath string, funcName string, params []interface{}, kw map[string]interface{}) PResult {
	result := PResult{}
	if _, err := os.Stat(scriptPath); os.IsNotExist(err) {
		result.Exception = fmt.Errorf("invoke python function error: python script not exists: %v: %v", scriptPath, err)
		return result
	}

	tempScript, appendPythonPath, err := generateTempScript(scriptPath, funcName, params, kw)
	if err != nil {
		result.Exception = fmt.Errorf("invoke python function error: generate temp script error: %v", err)
		return result
	}
	result.TempScript = tempScript

	cmd := exec.Command(PythonExecutable)

	cmd.Env = os.Environ()
	if len(appendPythonPath) > 0 {
		AddEnv(&cmd.Env, PythonPath, appendPythonPath)
	}
	result.PythonPath, _ = GetEnv(&cmd.Env, PythonPath)

	sin, err := cmd.StdinPipe()
	if err != nil {
		result.Exception = fmt.Errorf("invoke python function error: pipe stdin of python error: %v", err)
		return result
	}

	sout, err := cmd.StdoutPipe()
	if err != nil {
		result.Exception = fmt.Errorf("invoke python function error: pipe stdout of python error: %v", err)
		return result
	}
	defer sout.Close()

	serr, err := cmd.StderrPipe()
	if err != nil {
		result.Exception = fmt.Errorf("invoke python function error: pipe stderr of python error: %v", err)
		return result
	}
	defer serr.Close()

	err = cmd.Start()
	if err != nil {
		result.Exception = fmt.Errorf("invoke python function error: %v", err)
		return result
	}

	_, err = sin.Write([]byte(tempScript))
	if err != nil {
		result.Exception = fmt.Errorf("invoke python function error: input python script to python error: %v", err)
		return result
	}
	sin.Close()

	bs, err := ioutil.ReadAll(sout)
	if err != nil {
		result.Exception = fmt.Errorf("invoke python function error: get python output error: %v", err)
		return result
	}

	be, err := ioutil.ReadAll(serr)
	if err != nil {
		result.Exception = fmt.Errorf("invoke python function error: get python error output error: %v", err)
		return result
	}

	output := string(bs)
	errorOutput := string(be)

	result.Output = output + errorOutput
	result.JsonRepresentation = SubStringBetween(output, ReturnValueStart, ReturnValueEnd)
	result.Exception = errors.New(SubStringBetween(output, ExceptionStart, ExceptionEnd))

	if len(errorOutput) > 0 {
		result.Exception = errors.New(errorOutput)
	}

	if len(result.Exception.Error()) < 1 {
		result.NoError = true
	}

	return result
}

// generate temp script to send to python interpreter
func generateTempScript(scriptPath string, funcName string, params []interface{}, kw map[string]interface{}) (string, string, error) {
	script := bytes.Buffer{}
	var appendPythonPath string

	from, err := getRelativeImportPath(scriptPath)
	if err != nil {
		from, appendPythonPath = getAbsoluteImportPath(scriptPath)
	}

	vars, err := injectScriptVars(params, kw)
	if err != nil {
		return "", appendPythonPath, err
	}

	invoker, err := injectScriptFuncInvoke(funcName, params, kw)
	if err != nil {
		return "", appendPythonPath, err
	}

	str := fmt.Sprintf(Python2ScriptTemplate,
		from,
		funcName,
		TabString(vars, 4),
		invoker,
		ReturnValueStart,
		ReturnValueEnd,
		ExceptionStart,
		ExceptionEnd)

	script.WriteString(str)
	return script.String(), appendPythonPath, nil
}

func GetPythonPaths() []string {
	env := os.Environ()
	value, i := GetEnv(&env, PythonPath)
	if i < 0 || len(strings.TrimSpace(value)) < 1 {
		return []string{}
	} else {
		value = strings.TrimSpace(value)
		return strings.Split(value, string(os.PathListSeparator))
	}
}

// getRelativeImportPath get import path by relative path from PYTHONPATH to target script
func getRelativeImportPath(scriptPath string) (string, error) {
	ps := GetPythonPaths()
	for _, p := range ps {
		if b, rel := InPath(p, scriptPath); b {
			rel = strings.Replace(rel, string(os.PathSeparator)+"..", ".", -1)
			rel = strings.Replace(rel, string(os.PathSeparator), ".", -1)
			rel = strings.TrimSuffix(rel, ".py")
			//if !strings.HasSuffix(rel, "..") {
			//	rel = "." + rel
			//}
			return rel, nil
		}
	}
	return "", fmt.Errorf("cannot get relative import path by current PYTHONPATH")
}

// getAbsoluteImportPath get import path by add the dir of script to PYTHONPATH
func getAbsoluteImportPath(scriptPath string) (string, string) {
	p, _ := filepath.Abs(scriptPath)
	d, _ := filepath.Abs(filepath.Dir(p))
	b := filepath.Base(p)
	b = strings.TrimSuffix(b, ".py")
	return b, d
}

// injectScriptFuncInvoke generate script section to invoke and pass value to an python function,
// for example:
//   func1(var1, var2)
func injectScriptFuncInvoke(funcName string, params []interface{}, kw map[string]interface{}) (string, error) {
	var args []string
	for i, _ := range params {
		varName := fmt.Sprintf("%s%d", InjectVarNamePrefix, i)
		args = append(args, varName)
	}
	if kw != nil {
		for k, _ := range kw {
			varNameValue := fmt.Sprintf("%s = %s_%s", k, InjectVarNamePrefix, k)
			args = append(args, varNameValue)
		}
	}
	return fmt.Sprintf("%s(%s)", funcName, strings.Join(args, ", ")), nil
}

// injectScriptVars generate script section to define some variable. for example:
//   var1 = xxx
//   var2 = yyy
func injectScriptVars(params []interface{}, kw map[string]interface{}) (string, error) {
	if len(params) < 1 {
		return "", nil
	}

	script := bytes.Buffer{}

	for i, param := range params {
		varName := fmt.Sprintf("%s%d", InjectVarNamePrefix, i)
		bs, err := json.Marshal(param)
		if err != nil {
			return "", fmt.Errorf("can not serialize param to json value: %v", err)
		}
		varValue := string(bs)
		script.WriteString(fmt.Sprintf("%s = %s\n", varName, varValue))
	}

	if kw != nil {
		for k, v := range kw {
			varName := fmt.Sprintf("%s_%s", InjectVarNamePrefix, k)
			bs, err := json.Marshal(v)
			if err != nil {
				return "", fmt.Errorf("can not serialize keyword param to json value: %v", err)
			}
			varValue := string(bs)
			script.WriteString(fmt.Sprintf("%s = %s\n", varName, varValue))
		}
	}

	return script.String(), nil
}

// SubStringBetween return substring between inputted prefix and suffix
func SubStringBetween(str string, prefix string, suffix string) string {
	start := strings.Index(str, prefix)
	end := strings.Index(str, suffix)
	if start > -1 && end > -1 && start+len(prefix) < end {
		return str[start+len(prefix) : end]
	}
	return ""
}

// TabString function add some blank before every line in an string and return new string
func TabString(str string, tabSize int) string {
	result := bytes.NewBufferString("")

	buffer := bytes.NewBufferString(str)
	scanner := bufio.NewScanner(buffer)

	first := true
	for scanner.Scan() {
		if !first {
			result.WriteString("\n")
		} else {
			first = false
		}
		result.WriteString(strings.Repeat(" ", tabSize))
		result.WriteString(scanner.Text())
	}
	return result.String()
}

// Select function receive an bool flag as first parameter,
// and return the second or third parameter according to the flag,
// just like C style operator a ? b : c
func Select(tf bool, a interface{}, b interface{}) interface{} {
	if tf {
		return a
	} else {
		return b
	}
}

// FindLine return first line in string which contains inputted substring
func FindLine(str string, subs ...string) string {
	scanner := bufio.NewScanner(bytes.NewBufferString(str))
outer:
	for scanner.Scan() {
		for _, sub := range subs {
			if !strings.Contains(scanner.Text(), sub) {
				continue outer
			}
		}
		return scanner.Text()
	}
	return ""
}

// GetEnv get env value from env item slice
func GetEnv(env *[]string, key string) (string, int) {
	for i, item := range *env {
		if strings.HasPrefix(item, key+"=") {
			return strings.TrimPrefix(item, key+"="), i
		}
	}
	return "", -1
}

// AddEnv add env value to env item slice
func AddEnv(env *[]string, key string, value string) {
	oldValue, i := GetEnv(env, key)
	if len(oldValue) > 0 {
		(*env)[i] = (*env)[i] + string(os.PathListSeparator) + value
	} else if i > -1 {
		(*env)[i] = key + "=" + value
	} else {
		*env = append(*env, key+"="+value)
	}
}

// InPath check if the dir is parent of sub
func InPath(dir string, sub string) (bool, string) {
	dir, err := filepath.Abs(dir)
	if err != nil {
		return false, ""
	}
	sub, err = filepath.Abs(sub)
	if err != nil {
		return false, ""
	}
	rel, err := filepath.Rel(dir, sub)
	if err != nil {
		return false, ""
	}
	if !strings.HasPrefix(rel, ".") {
		return true, rel
	} else {
		return false, ""
	}
}
