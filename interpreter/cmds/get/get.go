package getcmd

import (
	"fmt"
	interpreter "gitlab.com/jbyte777/prompt-ql/v4/core"
)

func GetCmd(
	staticArgs interpreter.TFunctionArgumentsTable,
	_inputs interpreter.TFunctionInputChannelTable,
	internalGlobals interpreter.TGlobalVariablesTable,
	externalGlobals interpreter.TGlobalVariablesTable,
	execInfo interpreter.TExecutionInfo,
	_interpreter *interpreter.Interpreter,
) interface{} {
	fromVar, isExternal, err := getFromVar(staticArgs, execInfo)
	if err != nil {
		return err
	}

	varTable := internalGlobals
	if isExternal {
		varTable = externalGlobals
	}

	rawVar, hasVar := varTable[fromVar]
	if !hasVar {
		if isExternal {
			return fmt.Errorf(
				"!error (line=%v, char=%v): external variable with name \"%v\" is not defined",
				execInfo.Line,
				execInfo.CharPos,
				fromVar,
			)
		}
		
		return nil
	}
	return rawVar
}
