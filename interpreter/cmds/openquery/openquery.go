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
	
		model, err := getModel(staticArgs, execInfo)
		if err != nil {
			return err
		}
		
		temperature, err := getTemperature(staticArgs, execInfo)
		if err != nil {
			return err
		}
	
		prompts, err := getPrompts(inputs, execInfo)
		if err != nil {
			return err
		}

		queryHandle := gptApi.OpenQuery(
			model,
			temperature,
			prompts,
		)
		globals[toVar] = queryHandle

		return nil
	}
}
