package openquerycmd

import (
	api "gitlab.com/jbyte777/prompt-ql/api"
	interpreter "gitlab.com/jbyte777/prompt-ql/core"
	customapis "gitlab.com/jbyte777/prompt-ql/custom-apis"
)

func MakeOpenQueryCmd(
	gptApi *api.GptApi,
	customApis *customapis.CustomLLMApis,
) interpreter.TExecutedFunction {
	standardOpenQuery := func(
		model string,
		temperature float64,
		inputs interpreter.TFunctionInputChannelTable,
		execInfo interpreter.TExecutionInfo,
	) (*api.TQueryHandle, error) {
		prompts, err := getPrompts(inputs, execInfo)
		if err != nil {
			return nil, err
		}

		queryHandle := gptApi.OpenQuery(
			model,
			temperature,
			prompts,
		)
		return queryHandle, nil
	}

	userOpenQuery := func(
		model string,
		temperature float64,
		inputs interpreter.TFunctionInputChannelTable,
		execInfo interpreter.TExecutionInfo,
	) (*customapis.TCustomQueryHandle, error) {
		queryHandle, err := customApis.OpenQuery(
			model,
			temperature,
			inputs,
			execInfo,
		)
		return queryHandle, err
	}

	return func(
		staticArgs interpreter.TFunctionArgumentsTable,
		inputs interpreter.TFunctionInputChannelTable,
		globals interpreter.TGlobalVariablesTable,
		execInfo interpreter.TExecutionInfo,
	) interface{} {
		toVar, err := getToVar(staticArgs, execInfo)
		if err != nil {
			return err
		}

		userFlag := getUserFlag(staticArgs)
	
		model, err := getModel(staticArgs, execInfo)
		if err != nil {
			return err
		}
		
		temperature, err := getTemperature(staticArgs, execInfo)
		if err != nil {
			return err
		}

		if userFlag {
			queryHandle, err := userOpenQuery(
				model,
				temperature,
				inputs,
				execInfo,
			)
			if err != nil {
				return err
			}
			globals[toVar] = queryHandle
		} else {
			queryHandle, err := standardOpenQuery(
				model,
				temperature,
				inputs,
				execInfo,
			)
			if err != nil {
				return err
			}
			globals[toVar] = queryHandle
		}

		return nil
	}
}
