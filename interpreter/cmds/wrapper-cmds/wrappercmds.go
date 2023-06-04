package wrappercmds

import (
	"fmt"

	interpreter "gitlab.com/jbyte777/prompt-ql/core"
	promptmsg "gitlab.com/jbyte777/prompt-ql/utils/promptmsg"
)

func MakeWrapperCmd(dataTag string) interpreter.TExecutedFunction {
	return func(
		globals interpreter.TGlobalVariablesTable,
		staticArgs interpreter.TFunctionArgumentsTable,
		inputs interpreter.TFunctionInputChannelTable,
	) interface{} {
		dataChan, hasDataChan := inputs["data"]
		if !hasDataChan || len(dataChan) == 0 {
			return fmt.Errorf(
				"!error data is not provided for the \"%v\" wrapper",
				dataTag,
			)
		}
	
		rawLatestData := dataChan[len(dataChan) - 1]
		latestData, isLatestDataStr := rawLatestData.(string)
		if !isLatestDataStr {
			return fmt.Errorf(
				"!error \"%v\" is not valid string for \"%v\" wrapper",
				rawLatestData,
				dataTag,
			)
		}
		return promptmsg.ReplacePromptMsgPrefix(
			latestData,
			dataTag,
		)
	}
}
