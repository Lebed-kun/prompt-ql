package setcmd

import (
	"fmt"

	interpreter "gitlab.com/jbyte777/prompt-ql/core"
)

func SetCmd(
	staticArgs interpreter.TFunctionArgumentsTable,
	inputs interpreter.TFunctionInputChannelTable,
	internalGlobals interpreter.TGlobalVariablesTable,
	externalGlobals interpreter.TGlobalVariablesTable,
	execInfo interpreter.TExecutionInfo,
) interface{} {
	toVar, err := getToVar(staticArgs, execInfo)
	if err != nil {
		return err
	}

	dataChan, hasDataChan := inputs["data"]
	if !hasDataChan || len(dataChan) == 0 {
		return fmt.Errorf(
			"!error (line=%v, char=%v): data is not provided for the \"%v\" variable",
			execInfo.Line,
			execInfo.CharPos,
			toVar,
		)
	}

	latestData := dataChan[len(dataChan) - 1]
	internalGlobals[toVar] = latestData
	return nil
}
