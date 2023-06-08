package defaultswicth

import (
	interpreter "gitlab.com/jbyte777/prompt-ql/core"
	promptmsg "gitlab.com/jbyte777/prompt-ql/utils/promptmsg"
)

func DefaultSwitch(
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
		if isDataErr {
			topCtx.InputChannels["error"] = append(
				topCtx.InputChannels["error"],
				err.Error(),
			)
			return
		}

		topCtx.InputChannels["data"] = append(
			topCtx.InputChannels["data"],
			err,
		)
	} else {
		msgPrefix := promptmsg.GetPromptMsgPrefix(data)
		if len(msgPrefix) == 0 {
			msgPrefix = "data"
		}

		topCtx.InputChannels[msgPrefix] = append(
			topCtx.InputChannels[msgPrefix],
			promptmsg.ReplacePromptMsgPrefix(data, ""),
		)
	}
}
