package wrapperswicth

import (
	interpreter "gitlab.com/jbyte777/prompt-ql/v4/core"
	promptmsg "gitlab.com/jbyte777/prompt-ql/v4/utils/promptmsg"
	stringsutils "gitlab.com/jbyte777/prompt-ql/v4/utils/strings"
)

func WrapperSwitch(
	topCtx *interpreter.TExecutionStackFrame,
	rawData interface{},
) {
	_, hasMsgChan := topCtx.InputChannels["data"]
	if !hasMsgChan {
		topCtx.InputChannels["data"] = make(interpreter.TFunctionInputChannel, 0)
	}

	data, isDataStr := rawData.(string)
	if !isDataStr {
		err, isDataErr := rawData.(error)
		if isDataErr {
			topCtx.InputChannels["data"] = append(
				topCtx.InputChannels["data"],
				err.Error(),
			)
			return
		}

		topCtx.InputChannels["data"] = append(
			topCtx.InputChannels["data"],
			err,
		)
	} else {
		text := stringsutils.TrimWhitespace(
			promptmsg.ReplacePromptMsgPrefix(data, ""),
		)

		if len(text) > 0 && text != " " {
			topCtx.InputChannels["data"] = append(
				topCtx.InputChannels["data"],
				text,
			)
		}
	}
}
