package hellocmd

import (
	"encoding/json"
	"fmt"
	interpreter "gitlab.com/jbyte777/prompt-ql/v2/core"
	api "gitlab.com/jbyte777/prompt-ql/v2/api"
	customapis "gitlab.com/jbyte777/prompt-ql/v2/custom-apis"
)

func MakeHelloCmd(
	gptApi *api.GptApi,
	customApis *customapis.CustomModelsApis,
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
			MyModels: modelsList,
			MyVariables: externalGlobals,
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