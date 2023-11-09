package callcmd

import (
	"fmt"

	interpreter "gitlab.com/jbyte777/prompt-ql/v4/core"
)

func getFnVar(
	staticArgs interpreter.TFunctionArgumentsTable,
	execInfo interpreter.TExecutionInfo,
) (string, bool, error) {
	var fnVar string
	rawFnVar, hasRawFnVar := staticArgs["fn"]
	if !hasRawFnVar {
		return "", false, fmt.Errorf(
			"!error (line=%v, char=%v): \"fn\" parameter is required",
			execInfo.Line,
			execInfo.CharPos,
		)
	}
	var isFnVarStr bool
	fnVar, isFnVarStr = rawFnVar.(string)
	if !isFnVarStr || len(fnVar) == 0 {
		return "", false, fmt.Errorf(
			"!error (line=%v, char=%v): \"fn\" parameter is \"%v\" which is not a valid string",
			execInfo.Line,
			execInfo.CharPos,
			rawFnVar,
		)
	}

	isExternal := false
	if fnVar[0] == '@' {
		isExternal = true
		fnVar = fnVar[1:]
	}

	if len(fnVar) == 0 {
		return "", false, fmt.Errorf(
			"!error (line=%v, char=%v): \"fn\" parameter is \"%v\" which is not a valid string",
			execInfo.Line,
			execInfo.CharPos,
			rawFnVar,
		)
	}

	return fnVar, isExternal, nil
}
