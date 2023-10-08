package wrappercmds

import (
	"fmt"

	interpreter "gitlab.com/jbyte777/prompt-ql/v2/core"
	promptmsg "gitlab.com/jbyte777/prompt-ql/v2/utils/promptmsg"
)

func MakeWrapperCmd(dataTag string) interpreter.TExecutedFunction {
	return func(
		_staticArgs interpreter.TFunctionArgumentsTable,
		inputs interpreter.TFunctionInputChannelTable,
		_internalGlobals interpreter.TGlobalVariablesTable,
		_externalGlobals interpreter.TGlobalVariablesTable,
		execInfo interpreter.TExecutionInfo,
		_interpreter *interpreter.Interpreter,
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
