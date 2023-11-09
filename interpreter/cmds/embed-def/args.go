package embeddefcmd

import (
	"fmt"

	interpreter "gitlab.com/jbyte777/prompt-ql/v4/core"
)

func getNameVar(
	staticArgs interpreter.TFunctionArgumentsTable,
	execInfo interpreter.TExecutionInfo,
) (string, error) {
	var nameVar string
	rawNameVar, hasRawNameVar := staticArgs["name"]
	if !hasRawNameVar {
		return "", fmt.Errorf(
			"!error (line=%v, char=%v): \"name\" parameter is required",
			execInfo.Line,
			execInfo.CharPos,
		)
	}
	var isNameVarStr bool
	nameVar, isNameVarStr = rawNameVar.(string)
	if !isNameVarStr || len(nameVar) == 0 {
		return "", fmt.Errorf(
			"!error (line=%v, char=%v): \"name\" parameter is \"%v\" which is not a valid string",
			execInfo.Line,
			execInfo.CharPos,
			rawNameVar,
		)
	}

	return nameVar, nil
}

func getDescVar(
	staticArgs interpreter.TFunctionArgumentsTable,
	execInfo interpreter.TExecutionInfo,
) (string, error) {
	var descVar string
	rawDescVar, hasRawDescVar := staticArgs["desc"]
	if !hasRawDescVar {
		return "", nil
	}
	var isDescVarStr bool
	descVar, isDescVarStr = rawDescVar.(string)
	if !isDescVarStr {
		return "", fmt.Errorf(
			"!error (line=%v, char=%v): \"name\" parameter is \"%v\" which is not a valid string",
			execInfo.Line,
			execInfo.CharPos,
			rawDescVar,
		)
	}

	return descVar, nil
}
