package embedifcmd

import (
	"fmt"

	interpreter "gitlab.com/jbyte777/prompt-ql/v4/core"
)

func getCondVar(
	staticArgs interpreter.TFunctionArgumentsTable,
	execInfo interpreter.TExecutionInfo,
) (string, bool, error) {
	var condVar string
	rawCondVar, hasRawCondVar := staticArgs["cond"]
	if !hasRawCondVar {
		return "", false, fmt.Errorf(
			"!error (line=%v, char=%v): \"cond\" parameter is required",
			execInfo.Line,
			execInfo.CharPos,
		)
	}
	var isCondVarStr bool
	condVar, isCondVarStr = rawCondVar.(string)
	if !isCondVarStr || len(condVar) == 0 {
		return "", false, fmt.Errorf(
			"!error (line=%v, char=%v): \"cond\" parameter is \"%v\" which is not a valid string",
			execInfo.Line,
			execInfo.CharPos,
			rawCondVar,
		)
	}

	isExternal := false
	if condVar[0] == '@' {
		isExternal = true
		condVar = condVar[1:]
	}

	if len(condVar) == 0 {
		return "", false, fmt.Errorf(
			"!error (line=%v, char=%v): \"cond\" parameter is \"%v\" which is not a valid string",
			execInfo.Line,
			execInfo.CharPos,
			rawCondVar,
		)
	}

	return condVar, isExternal, nil
}
