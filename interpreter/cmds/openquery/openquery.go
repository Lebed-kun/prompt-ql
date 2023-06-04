package openquerycmd

import (
	api "gitlab.com/jbyte777/prompt-ql/api"
	interpreter "gitlab.com/jbyte777/prompt-ql/core"
)

func MakeOpenQueryCmd(
	gptApi *api.GptApi,
) interpreter.TExecutedFunction {
	return func(
		globals interpreter.TGlobalVariablesTable,
		staticArgs interpreter.TFunctionArgumentsTable,
		inputs interpreter.TFunctionInputChannelTable,
	) interface{} {
		toVar, err := getToVar(staticArgs)
		if err != nil {
			return err
		}
	
		model, err := getModel(staticArgs)
		if err != nil {
			return err
		}
		
		temperature, err := getTemperature(staticArgs)
		if err != nil {
			return err
		}
	
		prompts, err := getPrompts(inputs)
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
