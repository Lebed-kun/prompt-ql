package openquerytticmd

import (
	"fmt"
	"strconv"

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

func getWidth(
	staticArgs interpreter.TFunctionArgumentsTable,
	execInfo interpreter.TExecutionInfo,
) (uint, error) {
	var width string
	rawWidth, hasWidth := staticArgs["width"]
	if !hasWidth {
		return 640, nil
	}

	var isWidthStr bool
	width, isWidthStr = rawWidth.(string)
	if !isWidthStr {
		return 0, fmt.Errorf(
			"!error (line=%v, char=%v): \"width\" should be a numeric string",
			execInfo.Line,
			execInfo.CharPos,
			rawWidth,
		)
	}

	widthNum, err := strconv.Atoi(width)
	if err != nil {
		return 0, fmt.Errorf(
			"!error (line=%v, char=%v): %v",
			execInfo.Line,
			execInfo.CharPos,
			err.Error(),
		)
	}

	return uint(widthNum), nil
}

func getHeight(
	staticArgs interpreter.TFunctionArgumentsTable,
	execInfo interpreter.TExecutionInfo,
) (uint, error) {
	var height string
	rawHeight, hasHeight := staticArgs["height"]
	if !hasHeight {
		return 640, nil
	}

	var isHeightStr bool
	height, isHeightStr = rawHeight.(string)
	if !isHeightStr {
		return 0, fmt.Errorf(
			"!error (line=%v, char=%v): \"height\" should be a numeric string",
			execInfo.Line,
			execInfo.CharPos,
			rawHeight,
		)
	}

	heightNum, err := strconv.Atoi(height)
	if err != nil {
		return 0, fmt.Errorf(
			"!error (line=%v, char=%v): %v",
			execInfo.Line,
			execInfo.CharPos,
			err.Error(),
		)
	}

	return uint(heightNum), nil
}

func getResponseFormat(
	staticArgs interpreter.TFunctionArgumentsTable,
	execInfo interpreter.TExecutionInfo,
) (string, error) {
	var responseFormat string
	rawResponseFormat, hasResponseFormat := staticArgs["responseFormat"]
	if !hasResponseFormat {
		responseFormat = "url"
	} else {
		var isResponseFormatStr bool
		responseFormat, isResponseFormatStr = rawResponseFormat.(string)
		if !isResponseFormatStr {
			return "", fmt.Errorf(
				"!error (line=%v, char=%v): \"%v\" is not valid responseFormat name",
				execInfo.Line,
				execInfo.CharPos,
				rawResponseFormat,
			)
		}
	}

	return responseFormat, nil
}

func getSyncFlag(
	staticArgs interpreter.TFunctionArgumentsTable,
) bool {
	_, hasSync := staticArgs["sync"]
	return hasSync
}
