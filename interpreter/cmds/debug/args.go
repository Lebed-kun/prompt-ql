package debugcmd

import (
	"fmt"
	interpreter "gitlab.com/jbyte777/prompt-ql/v5/core"
)

func getLoggerVar(
	staticArgs interpreter.TFunctionArgumentsTable,
	execInfo interpreter.TExecutionInfo,
) (string, error) {
	var loggerVar string
	rawLoggerVar, hasRawLoggerVar := staticArgs["logger"]
	if !hasRawLoggerVar {
		return "default", nil
	}
	loggerVar, isLoggerVarStr := rawLoggerVar.(string)
	if !isLoggerVarStr || len(loggerVar) == 0 {
		return "", fmt.Errorf(
			"!error (line=%v, char=%v): \"logger\" parameter is \"%v\" which is not valid logger name",
			execInfo.Line,
			execInfo.CharPos,
			rawLoggerVar,
		)
	}

	return loggerVar, nil
}
