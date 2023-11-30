package embedifcmd

import (
	"fmt"

	interpreter "gitlab.com/jbyte777/prompt-ql/v5/core"
)

func EmbedIfCmd(
	staticArgs interpreter.TFunctionArgumentsTable,
	inputs interpreter.TFunctionInputChannelTable,
	internalGlobals interpreter.TGlobalVariablesTable,
	externalGlobals interpreter.TGlobalVariablesTable,
	execInfo interpreter.TExecutionInfo,
	_interpreter *interpreter.Interpreter,
) interface{} {
	condVar, isExternal, err := getCondVar(staticArgs, execInfo)
	if err != nil {
		return err
	}

	varTable := internalGlobals
	if isExternal {
		varTable = externalGlobals
	}

	rawCond, hasCond := varTable[condVar]
	if !hasCond {
		if isExternal {
			return fmt.Errorf(
				"!error (line=%v, char=%v): function with name \"%v\" doesn't exist in external variables",
				execInfo.Line,
				execInfo.CharPos,
				condVar,
			)
		}

		return fmt.Errorf(
			"!error (line=%v, char=%v): function with name \"%v\" doesn't exist in internal variables",
			execInfo.Line,
			execInfo.CharPos,
			condVar,
		)
	}

	cond, isCond := rawCond.(func([]interface{}) bool)
	if !isCond {
		if isExternal {
			return fmt.Errorf(
				"!error (line=%v, char=%v): external variable \"%v\" doesn't contain function",
				execInfo.Line,
				execInfo.CharPos,
				condVar,
			)
		}

		return fmt.Errorf(
			"!error (line=%v, char=%v): internal variable \"%v\" doesn't contain function",
			execInfo.Line,
			execInfo.CharPos,
			condVar,
		)
	}
	
	condInputs, hasCondInputs := inputs["data"]
	if !hasCondInputs || len(condInputs) < 2 {
		return fmt.Errorf(
			"!error (line=%v, char=%v): \"yes\"/\"no\" branches not passed",
			execInfo.Line,
			execInfo.CharPos,
		)
	}

	condResult := false
	if len(condInputs) == 2 {
		condResult = cond([]interface{}{})
	} else {
		condResult = cond(condInputs[:len(condInputs) - 2])
	}

	if condResult {
		return condInputs[len(condInputs) - 2]
	}
	return condInputs[len(condInputs) - 1]
}
