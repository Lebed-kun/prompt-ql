package getcmd

import (
	interpreter "gitlab.com/jbyte777/prompt-ql/core"
)

func GetCmd(
	staticArgs interpreter.TFunctionArgumentsTable,
	inputs interpreter.TFunctionInputChannelTable,
	globals interpreter.TGlobalVariablesTable,
	execInfo interpreter.TExecutionInfo,
) interface{} {
	fromVar, err := getFromVar(staticArgs, execInfo)
	if err != nil {
		return err
	}

	rawVar, hasVar := globals[fromVar]
	if !hasVar {
		return nil
	}
	return rawVar
}
