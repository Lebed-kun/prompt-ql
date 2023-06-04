package setcmd

import (
	"errors"
	"fmt"

	interpreter "gitlab.com/jbyte777/prompt-ql/core"
)

func getToVar(
	staticArgs interpreter.TFunctionArgumentsTable,
) (string, error) {
	var toVar string
	rawToVar, hasRawToVar := staticArgs["to"]
	if !hasRawToVar {
		return "", errors.New(
			"!error \"to\" parameter is required",
		)
	}
	var isToVarStr bool
	toVar, isToVarStr = rawToVar.(string)
	if !isToVarStr {
		return "", fmt.Errorf(
			"!error \"to\" parameter is \"%v\" which is not string",
			rawToVar,
		)
	}

	return toVar, nil
}
