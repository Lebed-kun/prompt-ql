package listenquerycmd

import (
	"fmt"
	api "gitlab.com/jbyte777/prompt-ql/v3/api"
	interpreter "gitlab.com/jbyte777/prompt-ql/v3/core"
	customapis "gitlab.com/jbyte777/prompt-ql/v3/custom-apis"
	utils "gitlab.com/jbyte777/prompt-ql/v3/interpreter/utils"
)

func MakeListenQueryCmd(
	gptApi *api.GptApi,
	customApis *customapis.CustomModelsApis,
) interpreter.TExecutedFunction {
	standardListenQuery := func(
		queryHandle *api.TQueryHandle,
		execInfo interpreter.TExecutionInfo,
	) interface{} {
		gptResponse, err := gptApi.ListenQuery(queryHandle)
		if err != nil {
			return fmt.Errorf(
				"!error (line=%v, char=%v): %v",
				execInfo.Line,
				execInfo.CharPos,
				err.Error(),
			)
		}

		return utils.MergeGptApiChoices(gptResponse.Choices)
	}

	userListenQuery := func(
		queryHandle *customapis.TCustomQueryHandle,
		execInfo interpreter.TExecutionInfo,
	) interface{} {
		llmResponse, err := customApis.ListenQuery(queryHandle)
		if err != nil {
			return fmt.Errorf(
				"!error (line=%v, char=%v): %v",
				execInfo.Line,
				execInfo.CharPos,
				err.Error(),
			)
		}

		return fmt.Sprintf("!assistant %v", llmResponse)
	}

	return func(
		staticArgs interpreter.TFunctionArgumentsTable,
		_inputs interpreter.TFunctionInputChannelTable,
		internalGlobals interpreter.TGlobalVariablesTable,
		_externalGlobals interpreter.TGlobalVariablesTable,
		execInfo interpreter.TExecutionInfo,
		_interpreter *interpreter.Interpreter,
	) interface{} {
		fromVar, err := getFromVar(staticArgs, execInfo)
		if err != nil {
			return err
		}

		rawQueryHandle, hasQueryHandle := internalGlobals[fromVar]
		if !hasQueryHandle {
			return fmt.Errorf(
				"!error (line=%v, char=%v): query handle by name \"%v\" doesn't exist in internal variables",
				execInfo.Line,
				execInfo.CharPos,
				fromVar,
			)
		}

		customQueryHandle, isCustomQueryHandle := rawQueryHandle.(*customapis.TCustomQueryHandle)
		if isCustomQueryHandle {
			return userListenQuery(customQueryHandle, execInfo)
		}

		queryHandle, isQueryHandleValid := rawQueryHandle.(*api.TQueryHandle)
		if !isQueryHandleValid {
			return fmt.Errorf(
				"!error (line=%v, char=%v): query handle by name \"%v\" is not valid as it's %v",
				execInfo.Line,
				execInfo.CharPos,
				fromVar,
				queryHandle,
			)
		}
		return standardListenQuery(queryHandle, execInfo)
	}
}
