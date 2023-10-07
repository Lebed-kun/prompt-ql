package hellocmd

import (
	"encoding/json"
	"fmt"
	interpreter "gitlab.com/jbyte777/prompt-ql/core"
	api "gitlab.com/jbyte777/prompt-ql/api"
	customapis "gitlab.com/jbyte777/prompt-ql/custom-apis"
)

func MakeHelloCmd(
	gptApi *api.GptApi,
	customApis *customapis.CustomLLMApis,
) interpreter.TExecutedFunction {
	return func(
		_staticArgs interpreter.TFunctionArgumentsTable,
		_inputs interpreter.TFunctionInputChannelTable,
		_internalGlobals interpreter.TGlobalVariablesTable,
		_externalGlobals interpreter.TGlobalVariablesTable,
		execInfo interpreter.TExecutionInfo,
		interpreter *interpreter.Interpreter,
	) interface{} {
		modelsList := gptApi.GetAllModelsList()
		customModels := customApis.GetAllModelsList()

		for model := range customModels {
			modelsList[model] = true
		}

		externalGlobals := interpreter.GetExternalGlobalsList()

		result := THelloCmdResponse{
			Models: modelsList,
			Variables: externalGlobals,
		}
		rawResult, err := json.Marshal(result)
		if err != nil {
			return fmt.Errorf(
				"!error (line=%v, char=%v): %v",
				execInfo.Line,
				execInfo.CharPos,
				err.Error(),
			)
		}

		return string(rawResult)
	}
}
