package wrappercmds

import (
	"fmt"

	interpreter "gitlab.com/jbyte777/prompt-ql/v5/core"
	promptmsg "gitlab.com/jbyte777/prompt-ql/v5/utils/promptmsg"
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
	
		latestData, err := dataChan.MergeIntoString()
		if err != nil {
			return fmt.Errorf(
				"!error (line=%v, char=%v): %v",
				execInfo.Line,
				execInfo.CharPos,
				err.Error(),
			)
		}

		// [BEGIN] TODO: include this in the 1.x-4.x patches 
		promptMsgPrefix := promptmsg.GetPromptMsgPrefix(latestData)
		if len(promptMsgPrefix) == 0 {
			return fmt.Sprintf("!%v %v", dataTag, latestData)
		}
		// [END] TODO: include this in the 1.x-4.x patches

		return promptmsg.ReplacePromptMsgPrefix(
			latestData,
			fmt.Sprintf("!%v", dataTag),
		)
	}
}
