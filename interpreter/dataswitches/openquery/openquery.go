package openqueryswicth

import (
	"fmt"

	interpreter "gitlab.com/jbyte777/prompt-ql/core"
	promptmsg "gitlab.com/jbyte777/prompt-ql/utils/promptmsg"
)

func OpenQuerySwitch(
	topCtx *interpreter.TExecutionStackFrame,
	rawData interface{},
) {
	_, hasErrorChan := topCtx.InputChannels["error"]
	if !hasErrorChan {
		topCtx.InputChannels["error"] = make(interpreter.TFunctionInputChannel, 0)
	}

	data, isDataStr := rawData.(string)
	if !isDataStr {
		err, isDataErr := rawData.(error)
		if !isDataErr {
			topCtx.InputChannels["error"] = append(
				topCtx.InputChannels["error"],
				fmt.Errorf(
					"Unknown data type for open_query command, met: %v",
					rawData,
				),
			)
			return
		}

		topCtx.InputChannels["error"] = append(
			topCtx.InputChannels["error"],
			err.Error(),
		)
	} else {
		msgPrefix := promptmsg.GetPromptMsgPrefix(data)
		if len(msgPrefix) == 0 {
			msgPrefix = "user"
		}

		_, hasMsgChan := topCtx.InputChannels[msgPrefix]
		if !hasMsgChan {
			topCtx.InputChannels[msgPrefix] = make(interpreter.TFunctionInputChannel, 0)
		}
		topCtx.InputChannels[msgPrefix] = append(
			topCtx.InputChannels[msgPrefix],
			promptmsg.ReplacePromptMsgPrefix(data, ""),
		)
	}
}
