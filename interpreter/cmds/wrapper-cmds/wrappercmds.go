package wrappercmds

import (
	"fmt"

	interpreter "gitlab.com/jbyte777/prompt-ql/core"
	promptmsg "gitlab.com/jbyte777/prompt-ql/utils/promptmsg"
)

func MakeWrapperCmd(dataTag string) interpreter.TExecutedFunction {
	return func(
		staticArgs interpreter.TFunctionArgumentsTable,
		inputs interpreter.TFunctionInputChannelTable,
		internalGlobals interpreter.TGlobalVariablesTable,
		externalGlobals interpreter.TGlobalVariablesTable,
		execInfo interpreter.TExecutionInfo,
	) interface{} {
		dataChan, hasDataChan := inputs["data"]
		if !hasDataChan || len(dataChan) == 0 {
			return fmt.Errorf(
				"!error (line=%v, char=%v): data is not provided for the \"%v\" wrapper",
				execInfo.Line,
				execInfo.CharPos,
				dataTag,
			)
		}
	
		rawLatestData := dataChan[len(dataChan) - 1]
		latestData, isLatestDataStr := rawLatestData.(string)
		if !isLatestDataStr {
			return fmt.Errorf(
				"!error (line=%v, char=%v): \"%v\" is not valid string for \"%v\" wrapper",
				execInfo.Line,
				execInfo.CharPos,
				rawLatestData,
				dataTag,
			)
		}
		return promptmsg.ReplacePromptMsgPrefix(
			latestData,
			fmt.Sprintf("!%v", dataTag),
		)
	}
}
