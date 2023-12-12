package blobreadfromfilecmd

import (
	"fmt"

	interpreter "gitlab.com/jbyte777/prompt-ql/v5/core"
)

func getPathVar(
	staticArgs interpreter.TFunctionArgumentsTable,
	execInfo interpreter.TExecutionInfo,
) (string, error) {
	var pathVar string
	rawPathVar, hasRawPathVar := staticArgs["path"]
	if !hasRawPathVar {
		return "", fmt.Errorf(
			"!error (line=%v, char=%v): \"path\" parameter is required",
			execInfo.Line,
			execInfo.CharPos,
		)
	}
	var isPathVarStr bool
	pathVar, isPathVarStr = rawPathVar.(string)
	if !isPathVarStr || len(pathVar) == 0 {
		return "", fmt.Errorf(
			"!error (line=%v, char=%v): \"path\" parameter is \"%v\" which is not a valid string",
			execInfo.Line,
			execInfo.CharPos,
			rawPathVar,
		)
	}

	return pathVar, nil
}
