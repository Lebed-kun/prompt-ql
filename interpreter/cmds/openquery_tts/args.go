package openqueryttscmd

import (
	"fmt"
	interpreter "gitlab.com/jbyte777/prompt-ql/v5/core"
)

func getToVar(
	staticArgs interpreter.TFunctionArgumentsTable,
	execInfo interpreter.TExecutionInfo,
	isSyncQuery bool,
) (string, error) {
	var toVar string
	rawToVar, hasRawToVar := staticArgs["to"]
	if !hasRawToVar && isSyncQuery {
		return "", nil
	}
	if !hasRawToVar {
		return "", fmt.Errorf(
			"!error (line=%v, char=%v): \"to\" parameter is required",
			execInfo.Line,
			execInfo.CharPos,
		)
	}
	var isToVarStr bool
	toVar, isToVarStr = rawToVar.(string)
	if !isToVarStr || len(toVar) == 0 {
		return "", fmt.Errorf(
			"!error (line=%v, char=%v): \"to\" parameter is \"%v\" which is not a valid string",
			execInfo.Line,
			execInfo.CharPos,
			rawToVar,
		)
	}

	if toVar[0] == '@' {
		return "", fmt.Errorf(
			"!error (line=%v, char=%v): \"to\" parameter is \"%v\" which is a name of external variable. Query handles can't be stored in external variables",
			execInfo.Line,
			execInfo.CharPos,
			rawToVar,
		)
	}

	return toVar, nil
}

func getModel(
	staticArgs interpreter.TFunctionArgumentsTable,
	execInfo interpreter.TExecutionInfo,
) (string, error) {
	var model string
	rawModel, hasModel := staticArgs["model"]
	if !hasModel {
		model = "tts-1"
	} else {
		var isModelStr bool
		model, isModelStr = rawModel.(string)
		if !isModelStr {
			return "", fmt.Errorf(
				"!error (line=%v, char=%v): \"%v\" is not valid model name",
				execInfo.Line,
				execInfo.CharPos,
				rawModel,
			)
		}
	}

	return model, nil
}

func getVoice(
	staticArgs interpreter.TFunctionArgumentsTable,
	execInfo interpreter.TExecutionInfo,
) (string, error) {
	var voice string
	rawVoice, hasVoice := staticArgs["voice"]
	if !hasVoice {
		return "", fmt.Errorf(
			"!error (line=%v, char=%v): \"voice\" parameter is required",
			execInfo.Line,
			execInfo.CharPos,
		)
	} else {
		var isVoiceStr bool
		voice, isVoiceStr = rawVoice.(string)
		if !isVoiceStr {
			return "", fmt.Errorf(
				"!error (line=%v, char=%v): \"%v\" is not valid voice name",
				execInfo.Line,
				execInfo.CharPos,
				rawVoice,
			)
		}
	}

	return voice, nil
}

func getSyncFlag(
	staticArgs interpreter.TFunctionArgumentsTable,
) bool {
	_, hasSync := staticArgs["sync"]
	return hasSync
}
