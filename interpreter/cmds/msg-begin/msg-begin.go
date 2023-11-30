package msgbegincmd

import (
	interpreter "gitlab.com/jbyte777/prompt-ql/v5/core"
)

func MsgBeginCmd(
	_staticArgs interpreter.TFunctionArgumentsTable,
	_inputs interpreter.TFunctionInputChannelTable,
	_internalGlobals interpreter.TGlobalVariablesTable,
	_externalGlobals interpreter.TGlobalVariablesTable,
	_execInfo interpreter.TExecutionInfo,
	_interpreter *interpreter.Interpreter,
) interface{} {
	return "[MSG_BEGIN]"
}
