package callcmd

import (
	"fmt"

	interpreter "gitlab.com/jbyte777/prompt-ql/core"
)

func CallCmd(
	staticArgs interpreter.TFunctionArgumentsTable,
	inputs interpreter.TFunctionInputChannelTable,
	globals interpreter.TGlobalVariablesTable,
	execInfo interpreter.TExecutionInfo,
) interface{} {
	fnVar, err := getFnVar(staticArgs, execInfo)
	if err != nil {
		return err
	}

	rawFn, hasFn := globals[fnVar]
	if !hasFn {
		return fmt.Errorf(
			"!error (line=%v, char=%v): function with name \"%v\" doesn't exist",
			execInfo.Line,
			execInfo.CharPos,
			fnVar,
		)
	}

	fn, isFn := rawFn.(TCmdCallableFunction)
	if !isFn {
		return fmt.Errorf(
			"!error (line=%v, char=%v): variable \"%v\" doesn't contain function",
			execInfo.Line,
			execInfo.CharPos,
			fnVar,
		)
	}

	fnInputs, hasFnInputs := inputs["data"]
	if !hasFnInputs {
		fnInputs = make(interpreter.TFunctionInputChannel, 0)
	}
	fnResult := fn(fnInputs)

	return fnResult
}
