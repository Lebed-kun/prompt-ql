package listenquerycmd

import (
	"fmt"
	"strings"

	api "gitlab.com/jbyte777/prompt-ql/api"
	interpreter "gitlab.com/jbyte777/prompt-ql/core"
)

func mergeChoices(choices []api.TGptApiResponseChoice) string {
	strChoices := make([]string, 0)
	for _, choice := range choices {
		if choice.Message.Role != "assistant" {
			continue
		}

		strChoices = append(strChoices, choice.Message.Content)
	}

	result := fmt.Sprintf(
		"!assistant %v",
		strings.Join(strChoices, "\n=====\n"),
	)
	return result
}

func MakeListenQueryCmd(
	gptApi *api.GptApi,
) interpreter.TExecutedFunction {
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

		gptResponse, err := gptApi.ListenQuery(queryHandle)
		if err != nil {
			return fmt.Errorf(
				"!error (line=%v, char=%v): %v",
				execInfo.Line,
				execInfo.CharPos,
				err.Error(),
			)
		}

		return mergeChoices(gptResponse.Choices)
	}
}
