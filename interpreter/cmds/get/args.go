package getcmd

import (
	"fmt"

	interpreter "gitlab.com/jbyte777/prompt-ql/v4/core"
)

func getFromVar(
	staticArgs interpreter.TFunctionArgumentsTable,
	execInfo interpreter.TExecutionInfo,
) (string, bool, error) {
	var fromVar string
	rawFromVar, hasRawFromVar := staticArgs["from"]
	if !hasRawFromVar {
		return "", false, fmt.Errorf(
			"!error (line=%v, char=%v): \"from\" parameter is required",
			execInfo.Line,
			execInfo.CharPos,
		)
	}
	var isFromVarStr bool
	fromVar, isFromVarStr = rawFromVar.(string)
	if !isFromVarStr || len(fromVar) == 0 {
		return "", false, fmt.Errorf(
			"!error (line=%v, char=%v): \"from\" parameter is \"%v\" which is not valid string",
			execInfo.Line,
			execInfo.CharPos,
			rawFromVar,
		)
	}

	isExternal := false
	if fromVar[0] == '@' {
		isExternal = true
		fromVar = fromVar[1:]
	}

	if len(fromVar) == 0 {
		return "", false, fmt.Errorf(
			"!error (line=%v, char=%v): \"from\" parameter is \"%v\" which is not a valid string",
			execInfo.Line,
			execInfo.CharPos,
			rawFromVar,
		)
	}

	return fromVar, isExternal, nil
}
