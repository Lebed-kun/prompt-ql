package headercmd

import (
	"encoding/json"
	"fmt"
	interpreter "gitlab.com/jbyte777/prompt-ql/v4/core"
)

func HeaderCmd(
	staticArgs interpreter.TFunctionArgumentsTable,
	_inputs interpreter.TFunctionInputChannelTable,
	_internalGlobals interpreter.TGlobalVariablesTable,
	_externalGlobals interpreter.TGlobalVariablesTable,
	execInfo interpreter.TExecutionInfo,
	_interpreter *interpreter.Interpreter,
) interface{} {
	rawFromName, isFromString := staticArgs["from"]
	if !isFromString || len(rawFromName.(string)) == 0 {
		return fmt.Errorf(
			"!error (line=%v, char=%v): \"%v\" is not a valid name of sender agent",
			execInfo.Line,
			execInfo.CharPos,
			rawFromName,
		)
	}

	rawToName, isToString := staticArgs["to"]
	if !isToString || len(rawToName.(string)) == 0 {
		return fmt.Errorf(
			"!error (line=%v, char=%v): \"%v\" is not a valid name of receiver agent",
			execInfo.Line,
			execInfo.CharPos,
			rawToName,
		)
	}

	result := THeaderCmdResponse{
		FromAgent: rawFromName.(string),
		ToAgent: rawToName.(string),
	}
	rawResult, err := json.Marshal(result)
	if err != nil {
		return fmt.Errorf(
			"!error (line=%v, char=%v): %v",
			execInfo.Line,
			execInfo.CharPos,
			err.Error(),
		)
	}

	return fmt.Sprintf(
		"HEADER:%v\n",
		string(rawResult),
	)
}
