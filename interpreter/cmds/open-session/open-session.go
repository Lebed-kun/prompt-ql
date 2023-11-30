package opensessioncmd

import (
	"fmt"
	interpreter "gitlab.com/jbyte777/prompt-ql/v5/core"
)

func OpenSessionCmd(
	_staticArgs interpreter.TFunctionArgumentsTable,
	_inputs interpreter.TFunctionInputChannelTable,
	_internalGlobals interpreter.TGlobalVariablesTable,
	_externalGlobals interpreter.TGlobalVariablesTable,
	execInfo interpreter.TExecutionInfo,
	interpreter *interpreter.Interpreter,
) interface{} {
	isSessClosed := interpreter.IsSessionClosed()
	if !isSessClosed {
		return fmt.Errorf(
			"!error (line=%v, char=%v): PromptQL session is already opened",
			execInfo.Line,
			execInfo.CharPos,
		)
	}
	interpreter.OpenSession()
	
	return nil
}
