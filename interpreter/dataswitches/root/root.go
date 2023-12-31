package rootswicth

import (
	interpreter "gitlab.com/jbyte777/prompt-ql/core"
	promptmsg "gitlab.com/jbyte777/prompt-ql/utils/promptmsg"
)

func RootSwitch(
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
		text := promptmsg.ReplacePromptMsgPrefix(data, "")

		if len(text) > 0 && text != " " {
			topCtx.InputChannels["data"] = append(
				topCtx.InputChannels["data"],
				text,
			)
		}
	}
}
