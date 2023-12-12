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

func getMethodVar(
	staticArgs interpreter.TFunctionArgumentsTable,
	execInfo interpreter.TExecutionInfo,
) (string, error) {
	var methodVar string
	rawMethodVar, hasRawMethodVar := staticArgs["method"]
	if !hasRawMethodVar {
		return "GET", nil
	}
	var isMethodVarStr bool
	methodVar, isMethodVarStr = rawMethodVar.(string)
	if !isMethodVarStr || len(methodVar) == 0 {
		return "", fmt.Errorf(
			"!error (line=%v, char=%v): \"method\" parameter is \"%v\" which is not a valid string",
			execInfo.Line,
			execInfo.CharPos,
			rawMethodVar,
		)
	}

	return methodVar, nil
}
