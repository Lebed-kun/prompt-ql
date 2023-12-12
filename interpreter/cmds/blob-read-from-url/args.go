package blobreadfromurlcmd

import (
	"fmt"

	interpreter "gitlab.com/jbyte777/prompt-ql/v5/core"
)

func getUrlVar(
	staticArgs interpreter.TFunctionArgumentsTable,
	execInfo interpreter.TExecutionInfo,
) (string, error) {
	var urlVar string
	rawUrlVar, hasRawUrlVar := staticArgs["url"]
	if !hasRawUrlVar {
		return "", fmt.Errorf(
			"!error (line=%v, char=%v): \"url\" parameter is required",
			execInfo.Line,
			execInfo.CharPos,
		)
	}
	var isUrlVarStr bool
	urlVar, isUrlVarStr = rawUrlVar.(string)
	if !isUrlVarStr || len(urlVar) == 0 {
		return "", fmt.Errorf(
			"!error (line=%v, char=%v): \"url\" parameter is \"%v\" which is not a valid string",
			execInfo.Line,
			execInfo.CharPos,
			rawUrlVar,
		)
	}

	return urlVar, nil
}
