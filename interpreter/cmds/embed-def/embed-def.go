package embeddefcmd

import (
	"fmt"

	interpreter "gitlab.com/jbyte777/prompt-ql/v4/core"
)

func EmbedDefCmd(
	staticArgs interpreter.TFunctionArgumentsTable,
	inputs interpreter.TFunctionInputChannelTable,
	_internalGlobals interpreter.TGlobalVariablesTable,
	_externalGlobals interpreter.TGlobalVariablesTable,
	execInfo interpreter.TExecutionInfo,
	interpreter *interpreter.Interpreter,
) interface{} {
	nameVar, err := getNameVar(staticArgs, execInfo)
	if err != nil {
		return err
	}
	descVar, err := getDescVar(staticArgs, execInfo)
	if err != nil {
		return err
	}

	dataChan, hasDataChan := inputs["data"]
	if !hasDataChan {
		return fmt.Errorf(
			"!error (line=%v, char=%v): embedding %v is empty",
			execInfo.Line,
			execInfo.CharPos,
			nameVar,
		)
	}
	dataChanStr, err := dataChan.MergeIntoString()
	if err != nil {
		return fmt.Errorf(
			"!error (line=%v, char=%v): %v",
			execInfo.Line,
			execInfo.CharPos,
			err.Error(),
		)
	}

	interpreter.RegisterEmbedding(nameVar, dataChanStr, descVar)
	return nil
}
