package embedexpcmd

import (
	"fmt"

	interpretermod "gitlab.com/jbyte777/prompt-ql/v5/core"
)

func EmbedExpCmd(
	staticArgs interpretermod.TFunctionArgumentsTable,
	inputs interpretermod.TFunctionInputChannelTable,
	_internalGlobals interpretermod.TGlobalVariablesTable,
	_externalGlobals interpretermod.TGlobalVariablesTable,
	execInfo interpretermod.TExecutionInfo,
	interpreter *interpretermod.Interpreter,
) interface{} {
	inlineFlag := getInlineFlag(staticArgs)

	nameVar, err := getNameVar(staticArgs, execInfo, inlineFlag)
	if err != nil {
		return err
	}

	embdInputs, hasEmbdInputs := inputs["data"]
	if !hasEmbdInputs || len(embdInputs) == 0 {
		if inlineFlag {
			return fmt.Errorf(
				"!error (line=%v, char=%v): inline embedding must have an embeddable code at least in inputs",
				execInfo.Line,
				execInfo.CharPos,
			)
		}

		embdInputs = make(interpretermod.TFunctionInputChannel, 0)
	}
	if inlineFlag {
		embdInputs = embdInputs[1:]
	}

	embdArgs, err := getArgs(embdInputs, execInfo)
	if err != nil {
		return err
	}

	var embdExp string
	if inlineFlag {
		embd, isEmbdStr := inputs["data"][0].(string)
		if !isEmbdStr {
			return fmt.Errorf(
				"!error (line=%v, char=%v): not valid code embedding: \"%v\"",
				execInfo.Line,
				execInfo.CharPos,
				inputs["data"][0],
			)
		}
		embdExp = interpreter.ExpandInlineEmbedding(embd, embdArgs)
	} else {
		embdExp, err = interpreter.ExpandEmbedding(nameVar, embdArgs)
		if err != nil {
			return fmt.Errorf(
				"!error (line=%v, char=%v): %v",
				execInfo.Line,
				execInfo.CharPos,
				err.Error(),
			)
		}
	}

	return fmt.Sprintf(
		"!data %v",
		embdExp,
	)
}
