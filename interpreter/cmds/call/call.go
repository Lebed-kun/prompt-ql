package callcmd

import (
	"fmt"

	interpreter "gitlab.com/jbyte777/prompt-ql/core"
)

func CallCmd(
	globals interpreter.TGlobalVariablesTable,
	staticArgs interpreter.TFunctionArgumentsTable,
	inputs interpreter.TFunctionInputChannelTable,
) interface{} {
	fnVar, err := getFnVar(staticArgs)
	if err != nil {
		return err
	}

	rawFn, hasFn := globals[fnVar]
	if !hasFn {
		return fmt.Errorf(
			"!error function with name \"%v\" doesn't exist",
			fnVar,
		)
	}

	fn, isFn := rawFn.(TCmdCallableFunction)
	if !isFn {
		return fmt.Errorf(
			"!error variable \"%v\" doesn't contain function",
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
