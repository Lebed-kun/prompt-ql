package debugswicth

import (
	interpreter "gitlab.com/jbyte777/prompt-ql/v5/core"
	promptmsg "gitlab.com/jbyte777/prompt-ql/v5/utils/promptmsg"
	stringsutils "gitlab.com/jbyte777/prompt-ql/v5/utils/strings"
)

var possibleChannels []string = []string{
	"user", "system", "assistant", "data", "error",
}

func pushDataToChan(
	chanTable interpreter.TFunctionInputChannelTable,
	chanName string,
	data interface{},
) {
	for _, ch := range possibleChannels {
		if ch == chanName {
			chanTable[ch] = append(chanTable[ch], data)
		} else {
			chanTable[ch] = append(chanTable[ch], "")
		}
	}
}

func DebugSwitch(
	topCtx *interpreter.TExecutionStackFrame,
	rawData interface{},
) {
	data, isDataStr := rawData.(string)
	if !isDataStr {
		err, isDataErr := rawData.(error)
		if !isDataErr {
			pushDataToChan(
				topCtx.InputChannels,
				"data",
				err,
			)
			return
		}

		_, hasErrorChan := topCtx.InputChannels["error"]
		if !hasErrorChan {
			topCtx.InputChannels["error"] = make(interpreter.TFunctionInputChannel, 0)
		}
		pushDataToChan(
			topCtx.InputChannels,
			"error",
			err.Error(),
		)
	} else {
		msgPrefix := promptmsg.GetPromptMsgPrefix(data)
		if len(msgPrefix) == 0 {
			msgPrefix = "data"
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

		_, hasDataChan := topCtx.InputChannels["data"]
		if !hasDataChan {
			topCtx.InputChannels["data"] = make(interpreter.TFunctionInputChannel, 0)
		}

		_, hasErrorChan := topCtx.InputChannels["error"]
		if !hasErrorChan {
			topCtx.InputChannels["error"] = make(interpreter.TFunctionInputChannel, 0)
		}

		message := stringsutils.TrimWhitespace(
			promptmsg.ReplacePromptMsgPrefix(data, ""),
		)

		if msgPrefix == "user" {
			pushDataToChan(
				topCtx.InputChannels,
				"user",
				message,
			)
		} else if msgPrefix == "system" {
			pushDataToChan(
				topCtx.InputChannels,
				"system",
				message,
			)
		} else if msgPrefix == "assistant" {
			pushDataToChan(
				topCtx.InputChannels,
				"assistant",
				message,
			)
		} else if msgPrefix == "error" {
			pushDataToChan(
				topCtx.InputChannels,
				"error",
				message,
			)
		} else {
			pushDataToChan(
				topCtx.InputChannels,
				"data",
				message,
			)
		}
	}
}
