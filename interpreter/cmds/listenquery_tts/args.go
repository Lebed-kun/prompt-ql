package listenqueryttscmd

import (
	"fmt"

	interpreter "gitlab.com/jbyte777/prompt-ql/v5/core"
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
	if !isFromVarStr || len(fromVar) == 0 {
		return "", fmt.Errorf(
			"!error (line=%v, char=%v): \"from\" parameter is \"%v\" which is not a valid string",
			execInfo.Line,
			execInfo.CharPos,
			rawFromVar,
		)
	}

	if fromVar[0] == '@' {
		return "", fmt.Errorf(
			"!error (line=%v, char=%v): \"from\" parameter is \"%v\" which is a name of external variable. Query handles can't be stored in external variables",
			execInfo.Line,
			execInfo.CharPos,
			rawFromVar,
		)
	}

	return fromVar, nil
}
