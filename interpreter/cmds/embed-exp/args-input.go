package embedexpcmd

import (
	"fmt"

	interpretermod "gitlab.com/jbyte777/prompt-ql/v4/core"
)

func getArgs(
	inputs interpretermod.TFunctionInputChannel,
	execInfo interpretermod.TExecutionInfo,
) (interpretermod.TEmbeddingArgsTable, error) {
	res := make(interpretermod.TEmbeddingArgsTable, 0)

	for _, inp := range inputs {
		inpStr, isInpStr := inp.(string)
		if !isInpStr {
			return nil, fmt.Errorf(
				"!error (line=%v, char=%v): argument pair \"%v\" is not a string",
				execInfo.Line,
				execInfo.CharPos,
				inp,
			)
		}

		ptr := 0
		beginArgName := ptr
		for ptr < len(inpStr) && inpStr[ptr] != '=' {
			ptr++
		}
		if ptr >= len(inpStr) - 1 {
			return nil, fmt.Errorf(
				"!error (line=%v, char=%v): argument pair \"%v\" does not contain argument value",
				execInfo.Line,
				execInfo.CharPos,
				inp,
			)
		}
		argName := inpStr[beginArgName:ptr]

		ptr++
		argVal := inpStr[ptr:]

		res[argName] = argVal
	}

	return res, nil
}
