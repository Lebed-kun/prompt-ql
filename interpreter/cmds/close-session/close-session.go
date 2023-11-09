package closesessioncmd

import (
	"fmt"
	interpreter "gitlab.com/jbyte777/prompt-ql/v4/core"
)

func CloseSessionCmd(
	_staticArgs interpreter.TFunctionArgumentsTable,
	_inputs interpreter.TFunctionInputChannelTable,
	_internalGlobals interpreter.TGlobalVariablesTable,
	_externalGlobals interpreter.TGlobalVariablesTable,
	execInfo interpreter.TExecutionInfo,
	interpreter *interpreter.Interpreter,
) interface{} {
	isSessClosed := interpreter.IsSessionClosed()
	if isSessClosed {
		return fmt.Errorf(
			"!error (line=%v, char=%v): PromptQL session is already closed",
			execInfo.Line,
			execInfo.CharPos,
		)
	}
	interpreter.CloseSession()
	
	return nil
}
