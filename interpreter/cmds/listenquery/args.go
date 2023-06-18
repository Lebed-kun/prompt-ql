package listenquerycmd

import (
	"fmt"

	interpreter "gitlab.com/jbyte777/prompt-ql/core"
)

func getFromVar(
	staticArgs interpreter.TFunctionArgumentsTable,
	execInfo interpreter.TExecutionInfo,
) (string, error) {
	var fromVar string
	rawFromVar, hasRawFromVar := staticArgs["from"]
	if !hasRawFromVar {
		return "", fmt.Errorf(
			"!error (line=%v, char=%v): \"from\" parameter is required",
			execInfo.Line,
			execInfo.CharPos,
		)
	}
	var isFromVarStr bool
	fromVar, isFromVarStr = rawFromVar.(string)
	if !isFromVarStr {
		return "", fmt.Errorf(
			"!error  (line=%v, char=%v): \"from\" parameter is \"%v\" which is not string",
			execInfo.Line,
			execInfo.CharPos,
			rawFromVar,
		)
	}

	return fromVar, nil
}
