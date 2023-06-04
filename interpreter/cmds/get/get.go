package getcmd

import (
	interpreter "gitlab.com/jbyte777/prompt-ql/core"
)

func GetCmd(
	globals interpreter.TGlobalVariablesTable,
	staticArgs interpreter.TFunctionArgumentsTable,
	inputs interpreter.TFunctionInputChannelTable,
) interface{} {
	fromVar, err := getFromVar(staticArgs)
	if err != nil {
		return err
	}

	rawVar, hasVar := globals[fromVar]
	if !hasVar {
		return nil
	}
	return rawVar
}
