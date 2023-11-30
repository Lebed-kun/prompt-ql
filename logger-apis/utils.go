package loggerapis

import (
	"fmt"

	interpreter "gitlab.com/jbyte777/prompt-ql/v5/core"
)

func defaultLogger(
	execInfo interpreter.TExecutionInfo,
	inputs interpreter.TFunctionInputChannelTable,
	internalGlobals interpreter.TGlobalVariablesTable,
	externalGlobals interpreter.TGlobalVariablesTable,
) error {
	execInfoStr := fmt.Sprintf(`
		=======-----=======
		Interpreter cursor:
			+ Line: %v
			+ Column: %v
			+ Character at: %v
	`,
		execInfo.Line,
		execInfo.CharPos,
		execInfo.StrPos,
	)
	inputsStr := fmt.Sprintf(`
		=======-----=======
		Input channels:
			+ USER: %v
			+ SYSTEM: %v
			+ ASSISTANT: %v
			+ DATA: %v
			+ ERROR: %v
	`,
		inputs["user"],
		inputs["system"],
		inputs["assistant"],
		inputs["data"],
		inputs["error"],
	)
	internalGlobalsStr := fmt.Sprintf(`
		=======-----=======
		Internal globals:
			+ %v
	`,
		internalGlobals,
	)
	externalGlobalsStr := fmt.Sprintf(`
		=======-----=======
		External globals:
			+ %v
	`,
		externalGlobals,
	)

	res := fmt.Sprintf(`
		=======^^^^^=======
		DEBUG INFO
%v
%v
%v
%v
		=======_____=======
	`,
	execInfoStr,
	inputsStr,
	internalGlobalsStr,
	externalGlobalsStr,
	)

	fmt.Print(res)
	return nil
}
