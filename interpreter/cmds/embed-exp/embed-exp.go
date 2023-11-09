package embedexpcmd

import (
	"fmt"

	interpretermod "gitlab.com/jbyte777/prompt-ql/v4/core"
)

func EmbedExpCmd(
	staticArgs interpretermod.TFunctionArgumentsTable,
	inputs interpretermod.TFunctionInputChannelTable,
	_internalGlobals interpretermod.TGlobalVariablesTable,
	_externalGlobals interpretermod.TGlobalVariablesTable,
	execInfo interpretermod.TExecutionInfo,
	interpreter *interpretermod.Interpreter,
) interface{} {
	nameVar, err := getNameVar(staticArgs, execInfo)
	if err != nil {
		return err
	}
	
	embdInputs, hasEmbdInputs := inputs["data"]
	if !hasEmbdInputs {
		embdInputs = make(interpretermod.TFunctionInputChannel, 0)
	}
	
	embdArgs, err := getArgs(embdInputs, execInfo)
	if err != nil {
		return err
	}

	embdExp, err := interpreter.ExpandEmbedding(nameVar, embdArgs)
	if err != nil {
		return fmt.Errorf(
			"!error (line=%v, char=%v): %v",
			execInfo.Line,
			execInfo.CharPos,
			err.Error(),
		)
	}

	return fmt.Sprintf(
		"!data %v",
		embdExp,
	)
}
