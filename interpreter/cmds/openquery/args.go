package openquerycmd

import (
	"fmt"
	"strconv"

	interpreter "gitlab.com/jbyte777/prompt-ql/core"
)

func getToVar(
	staticArgs interpreter.TFunctionArgumentsTable,
	execInfo interpreter.TExecutionInfo,
) (string, error) {
	var toVar string
	rawToVar, hasRawToVar := staticArgs["to"]
	if !hasRawToVar {
		return "", fmt.Errorf(
			"!error (line=%v, char=%v): \"to\" parameter is required",
			execInfo.Line,
			execInfo.CharPos,
		)
	}
	var isToVarStr bool
	toVar, isToVarStr = rawToVar.(string)
	if !isToVarStr {
		return "", fmt.Errorf(
			"!error (line=%v, char=%v): \"to\" parameter is \"%v\" which is not string",
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
		model = "gpt-3.5-turbo"
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

func getTemperature(
	staticArgs interpreter.TFunctionArgumentsTable,
	execInfo interpreter.TExecutionInfo,
) (float64, error) {
	var temperature float64
	rawTemperature, hasTemperature := staticArgs["temperature"]
	if !hasTemperature {
		temperature = 1.0
	} else {
		var isTemperatureStr bool
		temperatureStr, isTemperatureStr := rawTemperature.(string)
		if !isTemperatureStr {
			return 0.0, fmt.Errorf(
				"!error (line=%v, char=%v): \"%v\" is not valid temperature value",
				execInfo.Line,
				execInfo.CharPos,
				rawTemperature,
			)
		}

		var err error
		temperature, err = strconv.ParseFloat(temperatureStr, 64)
		if err != nil {
			return 0.0, fmt.Errorf(
				"!error (line=%v, char=%v): %v",
				execInfo.Line,
				execInfo.CharPos,
				err.Error(),
			)
		}
	}

	return temperature, nil
}
