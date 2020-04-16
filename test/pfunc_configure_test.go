package test

import (
	"github.com/gitpillow/pfunc"
	"github.com/stretchr/testify/assert"
	"testing"
)

func Test(t *testing.T) {
	assert.Equal(t, pfunc.InjectVarNamePrefixDefault, pfunc.GetInjectVarNamePrefix())
	assert.Equal(t, pfunc.ReturnValueStartDefault, pfunc.GetReturnValueStart())
	assert.Equal(t, pfunc.ReturnValueEndDefault, pfunc.GetReturnValueEnd())
	assert.Equal(t, pfunc.ExceptionStartDefault, pfunc.GetExceptionStart())
	assert.Equal(t, pfunc.ExceptionEndDefault, pfunc.GetExceptionEnd())

	pfunc.AddTemplateElementNamesPrefix("my_prefix_")

	assert.Equal(t, "my_prefix_"+pfunc.InjectVarNamePrefixDefault, pfunc.GetInjectVarNamePrefix())
	assert.Equal(t, "my_prefix_"+pfunc.ReturnValueStartDefault, pfunc.GetReturnValueStart())
	assert.Equal(t, "my_prefix_"+pfunc.ReturnValueEndDefault, pfunc.GetReturnValueEnd())
	assert.Equal(t, "my_prefix_"+pfunc.ExceptionStartDefault, pfunc.GetExceptionStart())
	assert.Equal(t, "my_prefix_"+pfunc.ExceptionEndDefault, pfunc.GetExceptionEnd())

	pfunc.ResetTemplateElementNames()

	assert.Equal(t, pfunc.InjectVarNamePrefixDefault, pfunc.GetInjectVarNamePrefix())
	assert.Equal(t, pfunc.ReturnValueStartDefault, pfunc.GetReturnValueStart())
	assert.Equal(t, pfunc.ReturnValueEndDefault, pfunc.GetReturnValueEnd())
	assert.Equal(t, pfunc.ExceptionStartDefault, pfunc.GetExceptionStart())
	assert.Equal(t, pfunc.ExceptionEndDefault, pfunc.GetExceptionEnd())

	pfunc.SetPythonExecutable("python27")
	assert.Equal(t, "python27", pfunc.GetPythonExecutable())

	pfunc.SetInjectVarNamePrefix("my_inject_var_name_prefix_")
	assert.Equal(t, "my_inject_var_name_prefix_", pfunc.GetInjectVarNamePrefix())

	pfunc.SetReturnValueStart("my_return_value_start_")
	assert.Equal(t, "my_return_value_start_", pfunc.GetReturnValueStart())

	pfunc.SetReturnValueEnd("my_return_value_end_")
	assert.Equal(t, "my_return_value_end_", pfunc.GetReturnValueEnd())

	pfunc.SetExceptionStart("my_exception_start_")
	assert.Equal(t, "my_exception_start_", pfunc.GetExceptionStart())

	pfunc.SetExceptionEnd("my_exception_end_")
	assert.Equal(t, "my_exception_end_", pfunc.GetExceptionEnd())

	pfunc.ResetTemplateElementNames()

	assert.Equal(t, pfunc.InjectVarNamePrefixDefault, pfunc.GetInjectVarNamePrefix())
	assert.Equal(t, pfunc.ReturnValueStartDefault, pfunc.GetReturnValueStart())
	assert.Equal(t, pfunc.ReturnValueEndDefault, pfunc.GetReturnValueEnd())
	assert.Equal(t, pfunc.ExceptionStartDefault, pfunc.GetExceptionStart())
	assert.Equal(t, pfunc.ExceptionEndDefault, pfunc.GetExceptionEnd())
}
