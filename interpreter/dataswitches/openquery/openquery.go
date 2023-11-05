package openqueryswicth

import (
	"fmt"

	interpreter "gitlab.com/jbyte777/prompt-ql/v3/core"
	promptmsg "gitlab.com/jbyte777/prompt-ql/v3/utils/promptmsg"
	stringsutils "gitlab.com/jbyte777/prompt-ql/v3/utils/strings"
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
		if len(msgPrefix) == 0 || msgPrefix == "data" {
			msgPrefix = "user"
		}

		_, hasUserChan := topCtx.InputChannels["user"]
		if !hasUserChan {
			topCtx.InputChannels["user"] = make(interpreter.TFunctionInputChannel, 0)
		}

		_, hasAssistantChan := topCtx.InputChannels["assistant"]
		if !hasAssistantChan {
			topCtx.InputChannels["assistant"] = make(interpreter.TFunctionInputChannel, 0)
		}

		_, hasSystemChan := topCtx.InputChannels["system"]
		if !hasSystemChan {
			topCtx.InputChannels["system"] = make(interpreter.TFunctionInputChannel, 0)
		}

		message := stringsutils.TrimWhitespace(
			promptmsg.ReplacePromptMsgPrefix(data, ""),
		)
		if len(message) == 0 {
			return
		}

		if msgPrefix == "user" {
			topCtx.InputChannels["user"] = append(topCtx.InputChannels["user"], message)
			topCtx.InputChannels["assistant"] = append(topCtx.InputChannels["assistant"], "")
			topCtx.InputChannels["system"] = append(topCtx.InputChannels["system"], "")
		} else if msgPrefix == "assistant" {
			topCtx.InputChannels["user"] = append(topCtx.InputChannels["user"], "")
			topCtx.InputChannels["assistant"] = append(topCtx.InputChannels["assistant"], message)
			topCtx.InputChannels["system"] = append(topCtx.InputChannels["system"], "")
		} else {
			topCtx.InputChannels["user"] = append(topCtx.InputChannels["user"], "")
			topCtx.InputChannels["assistant"] = append(topCtx.InputChannels["assistant"], "")
			topCtx.InputChannels["system"] = append(topCtx.InputChannels["system"], message)
		}
	}
}
