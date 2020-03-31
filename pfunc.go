package pfunc

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

const InjectVarNamePrefix string = "pfunc_inject_"
const ReturnValueStart string = "pfunc_return_start_"
const ReturnValueEnd string = "pfunc_return_end_"
const ExceptionStart string = "pfunc_exception_start_"
const ExceptionEnd string = "pfunc_exception_end_"
const PythonExecutable string = "python"

func Invoke(scriptPath string, funcName string, params []interface{}) (interface{}, string, error) {
	if _, err := os.Stat(scriptPath); os.IsNotExist(err) {
		return nil, "", fmt.Errorf("invoke python function error: python script not exists: %v: %v", scriptPath, err)
	}

	tempScript, err := generateTempScript(scriptPath, funcName, params)
	if err != nil {
		return nil, "", fmt.Errorf("invoke python function error: generate temp script error: %v", err)
	}

	cmd := exec.Command(PythonExecutable)

	sin, err := cmd.StdinPipe()
	if err != nil {
		return nil, tempScript, fmt.Errorf("invoke python function error: pipe stdin of python error: %v", err)
	}

	sout, err := cmd.StdoutPipe()
	if err != nil {
		return nil, tempScript, fmt.Errorf("invoke python function error: pipe stdout of python error: %v", err)
	}

	err = cmd.Start()
	if err != nil {
		return nil, tempScript, fmt.Errorf("invoke python function error: %v", err)
	}

	_, err = sin.Write([]byte(tempScript))
	if err != nil {
		return nil, tempScript, fmt.Errorf("invoke python function error: input python script to python error: %v", err)
	}
	sin.Close()

	bs, err := ioutil.ReadAll(sout)
	if err != nil {
		return nil, tempScript, fmt.Errorf("invoke python function error: get python output error: %v", err)
	}

	output := string(bs)
	result := SubStringBetween(output, ReturnValueStart, ReturnValueEnd)
	exception := SubStringBetween(output, ExceptionStart, ExceptionEnd)
	if len(exception) > 0 {
		return nil, tempScript, fmt.Errorf("%s", exception)
	} else if len(result) > 0 {
		var r interface{}
		err := json.Unmarshal([]byte(result), &r)
		if err != nil {
			return nil, tempScript, fmt.Errorf("invoke python function error: get python output error: %v", err)
		}
		return r, tempScript, nil
	} else {
		return nil, tempScript, nil
	}
}

func generateTempScript(scriptPath string, funcName string, params []interface{}) (string, error) {
	script := bytes.Buffer{}

	rel, err := getRelativeImportPath(scriptPath)
	if err != nil {
		return "", err
	}
	script.WriteString(fmt.Sprintf("from %s import %s\n", rel, funcName))

	vars, err := injectScriptVars(params)
	if err != nil {
		return "", err
	}
	script.WriteString(fmt.Sprintf("%s\n", vars))

	invoker, err := injectScriptFuncInvoke(funcName, params)
	if err != nil {
		return "", err
	}

	str := `
import traceback
import json
try:
	result = {1}
	print "{2}"
	print json.dumps(result)
	print "{3}"
except Exception, e:
	msg = traceback.format_exc()
	print "{4}"
	print msg
	print "{5}"
`
	str = strings.Replace(str, "{1}", invoker, 1)
	str = strings.Replace(str, "{2}", ReturnValueStart, 1)
	str = strings.Replace(str, "{3}", ReturnValueEnd, 1)
	str = strings.Replace(str, "{4}", ExceptionStart, 1)
	str = strings.Replace(str, "{5}", ExceptionEnd, 1)

	script.WriteString(str)
	return script.String(), nil
}

func getRelativeImportPath(scriptPath string) (string, error) {
	wd, err := os.Getwd()
	if err != nil {
		return "", fmt.Errorf("can not get relative import path of script %v: %v", scriptPath, err)
	}
	abs, _ := filepath.Abs(scriptPath)
	rel, err := filepath.Rel(wd, abs)
	if err != nil {
		return "", fmt.Errorf("can not get relative import path of script %v: %v", scriptPath, err)
	}
	rel = strings.Replace(rel, "/..", ".", -1)
	rel = strings.TrimSuffix(rel, ".py")
	return rel, nil
}

func injectScriptFuncInvoke(funcName string, params []interface{}) (string, error) {
	var args []string
	for i, _ := range params {
		varName := fmt.Sprintf("%s%d", InjectVarNamePrefix, i)
		args = append(args, varName)
	}
	return fmt.Sprintf("%s(%s)", funcName, strings.Join(args, ", ")), nil
}

func injectScriptVars(params []interface{}) (string, error) {
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

	return script.String(), nil
}

func SubStringBetween(str string, prefix string, suffix string) string {
	start := strings.Index(str, prefix)
	end := strings.Index(str, suffix)
	if start > -1 && end > -1 && start+len(prefix) < end {
		return str[start+len(prefix) : end]
	}
	return ""
}
