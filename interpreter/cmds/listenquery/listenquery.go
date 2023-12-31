package listenquerycmd

import (
	"fmt"
	api "gitlab.com/jbyte777/prompt-ql/api"
	interpreter "gitlab.com/jbyte777/prompt-ql/core"
	customapis "gitlab.com/jbyte777/prompt-ql/custom-apis"
	utils "gitlab.com/jbyte777/prompt-ql/interpreter/utils"
)

func MakeListenQueryCmd(
	gptApi *api.GptApi,
	customApis *customapis.CustomLLMApis,
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
		inputs interpreter.TFunctionInputChannelTable,
		globals interpreter.TGlobalVariablesTable,
		execInfo interpreter.TExecutionInfo,
	) interface{} {
		fromVar, err := getFromVar(staticArgs, execInfo)
		if err != nil {
			return err
		}

		rawQueryHandle, hasQueryHandle := globals[fromVar]
		if !hasQueryHandle {
			return fmt.Errorf(
				"!error (line=%v, char=%v): query handle by name \"%v\" doesn't exist",
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
