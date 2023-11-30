package setcmd

import (
	"fmt"

	interpreter "gitlab.com/jbyte777/prompt-ql/v5/core"
)

func getToVar(
	staticArgs interpreter.TFunctionArgumentsTable,
	execInfo interpreter.TExecutionInfo,
) (string, error) {
	var toVar string
	rawToVar, hasRawToVar := staticArgs["to"]
	if !hasRawToVar {
		return "", fmt.Errorf(
			"!error (line=%v, char=%v): \"to\" parameter is required",
			execInfo.Line,
			execInfo.CharPos,
		)
	}
	var isToVarStr bool
	toVar, isToVarStr = rawToVar.(string)
	if !isToVarStr || len(toVar) == 0 {
		return "", fmt.Errorf(
			"!error (line=%v, char=%v): \"to\" parameter is \"%v\" which is not string",
			execInfo.Line,
			execInfo.CharPos,
			rawToVar,
		)
	}

	if toVar[0] == '@' {
		return "", fmt.Errorf(
			"!error (line=%v, char=%v): \"to\" parameter is \"%v\" which is a name of external variable. Writes to external variables in PromptQL code are restricted",
			execInfo.Line,
			execInfo.CharPos,
			rawToVar,
		)
	}

	return toVar, nil
}
