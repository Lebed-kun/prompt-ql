package callcmd

import (
	"fmt"

	interpreter "gitlab.com/jbyte777/prompt-ql/core"
)

func getFnVar(
	staticArgs interpreter.TFunctionArgumentsTable,
	execInfo interpreter.TExecutionInfo,
) (string, error) {
	var fnVar string
	rawFnVar, hasRawFnVar := staticArgs["fn"]
	if !hasRawFnVar {
		return "", fmt.Errorf(
			"!error (line=%v, char=%v): \"fn\" parameter is required",
			execInfo.Line,
			execInfo.CharPos,
		)
	}
	var isFnVarStr bool
	fnVar, isFnVarStr = rawFnVar.(string)
	if !isFnVarStr {
		return "", fmt.Errorf(
			"!error (line=%v, char=%v): \"fn\" parameter is \"%v\" which is not string",
			execInfo.Line,
			execInfo.CharPos,
			rawFnVar,
		)
	}

	return fnVar, nil
}
