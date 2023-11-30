package clearvarscmd

import (
	interpreter "gitlab.com/jbyte777/prompt-ql/v5/core"
)

func UnsafeClearVarsCmd(
	_staticArgs interpreter.TFunctionArgumentsTable,
	_inputs interpreter.TFunctionInputChannelTable,
	_internalGlobals interpreter.TGlobalVariablesTable,
	_externalGlobals interpreter.TGlobalVariablesTable,
	_execInfo interpreter.TExecutionInfo,
	interpreter *interpreter.Interpreter,
) interface{} {
	interpreter.ControlFlowClearInternalVars()
	return nil
}
