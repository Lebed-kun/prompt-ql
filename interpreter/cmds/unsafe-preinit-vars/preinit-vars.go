package preinitvarscmd

import (
	interpreter "gitlab.com/jbyte777/prompt-ql/v5/core"
)

func UnsafePreinitVarsCmd(
	_staticArgs interpreter.TFunctionArgumentsTable,
	_inputs interpreter.TFunctionInputChannelTable,
	_internalGlobals interpreter.TGlobalVariablesTable,
	_externalGlobals interpreter.TGlobalVariablesTable,
	_execInfo interpreter.TExecutionInfo,
	interpreter *interpreter.Interpreter,
) interface{} {
	interpreter.ControlFlowPreinitInternalVars()
	return nil
}
