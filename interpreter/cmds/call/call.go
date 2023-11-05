package callcmd

import (
	"fmt"

	interpreter "gitlab.com/jbyte777/prompt-ql/v3/core"
)

func CallCmd(
	staticArgs interpreter.TFunctionArgumentsTable,
	inputs interpreter.TFunctionInputChannelTable,
	internalGlobals interpreter.TGlobalVariablesTable,
	externalGlobals interpreter.TGlobalVariablesTable,
	execInfo interpreter.TExecutionInfo,
	_interpreter *interpreter.Interpreter,
) interface{} {
	fnVar, isExternal, err := getFnVar(staticArgs, execInfo)
	if err != nil {
		return err
	}

	varTable := internalGlobals
	if isExternal {
		varTable = externalGlobals
	}

	rawFn, hasFn := varTable[fnVar]
	if !hasFn {
		if isExternal {
			return fmt.Errorf(
				"!error (line=%v, char=%v): function with name \"%v\" doesn't exist in external variables",
				execInfo.Line,
				execInfo.CharPos,
				fnVar,
			)
		}

		return fmt.Errorf(
			"!error (line=%v, char=%v): function with name \"%v\" doesn't exist in internal variables",
			execInfo.Line,
			execInfo.CharPos,
			fnVar,
		)
	}

	fn, isFn := rawFn.(func([]interface{}) interface{})
	if !isFn {
		if isExternal {
			return fmt.Errorf(
				"!error (line=%v, char=%v): external variable \"%v\" doesn't contain function",
				execInfo.Line,
				execInfo.CharPos,
				fnVar,
			)
		}

		return fmt.Errorf(
			"!error (line=%v, char=%v): internal variable \"%v\" doesn't contain function",
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
