package embedexpcmd

import (
	"fmt"

	interpreter "gitlab.com/jbyte777/prompt-ql/v4/core"
)

func getInlineFlag(staticArgs interpreter.TFunctionArgumentsTable) bool {
	_, hasInlineFlag := staticArgs["inline"]
	return hasInlineFlag
}

func getNameVar(
	staticArgs interpreter.TFunctionArgumentsTable,
	execInfo interpreter.TExecutionInfo,
	isInlineEmbed bool,
) (string, error) {
	var nameVar string
	rawNameVar, hasRawNameVar := staticArgs["name"]
	if !hasRawNameVar {
		if isInlineEmbed {
			return "", nil
		}

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
