package openquerycmd

import (
	"fmt"
	api "gitlab.com/jbyte777/prompt-ql/api"
	interpreter "gitlab.com/jbyte777/prompt-ql/core"
	customapis "gitlab.com/jbyte777/prompt-ql/custom-apis"
	utils "gitlab.com/jbyte777/prompt-ql/interpreter/utils"
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
		isSync bool,
	) (interface{}, error) {
		prompts, err := getPrompts(inputs, execInfo)
		if err != nil {
			return nil, err
		}

		if isSync {
			handle, err := gptApi.OpenQuery(
				model,
				temperature,
				prompts,
			)
			if err != nil {
				return nil, err
			}
			response, err := gptApi.ListenQuery(handle)
			if err != nil {
				return nil, err
			}
			return utils.MergeGptApiChoices(response.Choices), nil
		} else {
			return gptApi.OpenQuery(
				model,
				temperature,
				prompts,
			)
		}
	}

	userOpenQuery := func(
		model string,
		temperature float64,
		inputs interpreter.TFunctionInputChannelTable,
		execInfo interpreter.TExecutionInfo,
		isSync bool,
		) (interface{}, error) {
		if isSync {
			handle, err := customApis.OpenQuery(
				model,
				temperature,
				inputs,
				execInfo,
			)
			if err != nil {
				return nil, err
			}
			
			llmResponse, err := customApis.ListenQuery(handle)
			if err != nil {
				return nil, err
			}
			return fmt.Sprintf("!assistant %v", llmResponse), nil
		} else {
			return customApis.OpenQuery(
				model,
				temperature,
				inputs,
				execInfo,
			)
		}
	}

	return func(
		staticArgs interpreter.TFunctionArgumentsTable,
		inputs interpreter.TFunctionInputChannelTable,
		globals interpreter.TGlobalVariablesTable,
		execInfo interpreter.TExecutionInfo,
	) interface{} {
		userFlag := getUserFlag(staticArgs)
		syncFlag := getSyncFlag(staticArgs)

		toVar, err := getToVar(staticArgs, execInfo, syncFlag)
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

		isOpenAIModelSupported := gptApi.IsModelSupported(model)
		if userFlag || !isOpenAIModelSupported {
			queryHandleOrResponse, err := userOpenQuery(
				model,
				temperature,
				inputs,
				execInfo,
				syncFlag,
			)

			if err != nil {
				return fmt.Errorf(
					"!error (line=%v, char=%v): %v",
					execInfo.Line,
					execInfo.CharPos,
					err.Error(),
				)
			}
			
			if !syncFlag {
				globals[toVar] = queryHandleOrResponse
			} else {
				return queryHandleOrResponse
			}
		} else {
			queryHandleOrResponse, err := standardOpenQuery(
				model,
				temperature,
				inputs,
				execInfo,
				syncFlag,
			)

			if err != nil {
				return fmt.Errorf(
					"!error (line=%v, char=%v): %v",
					execInfo.Line,
					execInfo.CharPos,
					err.Error(),
				)
			}

			if !syncFlag {
				globals[toVar] = queryHandleOrResponse
			} else {
				return queryHandleOrResponse
			}
		}

		return nil
	}
}
