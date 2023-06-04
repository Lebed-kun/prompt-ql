package openquerycmd

import (
	"fmt"
	"strconv"
	"errors"

	interpreter "gitlab.com/jbyte777/prompt-ql/core"
)

func getToVar(
	staticArgs interpreter.TFunctionArgumentsTable,
) (string, error) {
	var toVar string
	rawToVar, hasRawToVar := staticArgs["to"]
	if !hasRawToVar {
		return "", errors.New(
			"!error \"to\" parameter is required",
		)
	}
	var isToVarStr bool
	toVar, isToVarStr = rawToVar.(string)
	if !isToVarStr {
		return "", fmt.Errorf(
			"!error \"to\" parameter is \"%v\" which is not string",
			rawToVar,
		)
	}

	return toVar, nil
}

func getModel(
	staticArgs interpreter.TFunctionArgumentsTable,
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
				"!error \"%v\" is not valid model name",
				rawModel,
			)
		}
	}

	return model, nil
}

func getTemperature(
	staticArgs interpreter.TFunctionArgumentsTable,
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
				"!error \"%v\" is not valid temperature value",
				rawTemperature,
			)
		}

		var err error
		temperature, err = strconv.ParseFloat(temperatureStr, 64)
		if err != nil {
			return 0.0, fmt.Errorf(
				"!error %v",
				err.Error(),
			)
		}
	}

	return temperature, nil
}
