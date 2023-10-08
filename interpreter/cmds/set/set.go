package setcmd

import (
	"fmt"

	interpreter "gitlab.com/jbyte777/prompt-ql/v2/core"
)

func SetCmd(
	staticArgs interpreter.TFunctionArgumentsTable,
	inputs interpreter.TFunctionInputChannelTable,
	internalGlobals interpreter.TGlobalVariablesTable,
	_externalGlobals interpreter.TGlobalVariablesTable,
	execInfo interpreter.TExecutionInfo,
	_interpreter *interpreter.Interpreter,
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
