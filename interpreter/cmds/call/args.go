package callcmd

import (
	"errors"
	"fmt"

	interpreter "gitlab.com/jbyte777/prompt-ql/core"
)

func getFnVar(
	staticArgs interpreter.TFunctionArgumentsTable,
) (string, error) {
	var fnVar string
	rawFnVar, hasRawFnVar := staticArgs["fn"]
	if !hasRawFnVar {
		return "", errors.New(
			"!error \"fn\" parameter is required",
		)
	}
	var isFnVarStr bool
	fnVar, isFnVarStr = rawFnVar.(string)
	if !isFnVarStr {
		return "", fmt.Errorf(
			"!error \"fn\" parameter is \"%v\" which is not string",
			rawFnVar,
		)
	}

	return fnVar, nil
}
