package setcmd

import (
	"fmt"

	interpreter "gitlab.com/jbyte777/prompt-ql/core"
)

func SetCmd(
	globals interpreter.TGlobalVariablesTable,
	staticArgs interpreter.TFunctionArgumentsTable,
	inputs interpreter.TFunctionInputChannelTable,
) interface{} {
	toVar, err := getToVar(staticArgs)
	if err != nil {
		return err
	}

	dataChan, hasDataChan := inputs["data"]
	if !hasDataChan || len(dataChan) == 0 {
		return fmt.Errorf(
			"!error data is not provided for the \"%v\" variable",
			toVar,
		)
	}

	latestData := dataChan[len(dataChan) - 1]
	globals[toVar] = latestData
	return nil
}
