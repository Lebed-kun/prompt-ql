package hellocmd

import (
	"encoding/json"
	"fmt"
	interpreter "gitlab.com/jbyte777/prompt-ql/v5/core"
	chatapi "gitlab.com/jbyte777/prompt-ql/v5/default-apis/chat-api"
	customapis "gitlab.com/jbyte777/prompt-ql/v5/custom-apis"
)

func MakeHelloCmd(
	gptApi *chatapi.GptApi,
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

		for name, desc := range customModels {
			modelsList[name] = desc
		}

		externalGlobals := interpreter.GetExternalGlobalsList()
		embeddings := interpreter.GetEmbeddingsList()

		result := THelloCmdResponse{
			MyModels: modelsList,
			MyVariables: externalGlobals,
			MyEmbeddings: embeddings,
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

		return fmt.Sprintf(
			"MY_LAYOUT:%v\n",
			string(rawResult),
		)
	}
}
