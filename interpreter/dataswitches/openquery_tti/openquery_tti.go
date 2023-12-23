package openqueryttiswicth

import (
	"fmt"

	interpreter "gitlab.com/jbyte777/prompt-ql/v5/core"
	promptmsg "gitlab.com/jbyte777/prompt-ql/v5/utils/promptmsg"
	stringsutils "gitlab.com/jbyte777/prompt-ql/v5/utils/strings"
)

func OpenQueryTtiSwitch(
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
					"Unknown data type for open_query_tti command, met: %v",
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
		_, hasDataChan := topCtx.InputChannels["data"]
		if !hasDataChan {
			topCtx.InputChannels["data"] = make(interpreter.TFunctionInputChannel, 0)
		}

		message := stringsutils.TrimWhitespace(
			promptmsg.ReplacePromptMsgPrefix(data, ""),
		)
		if len(message) == 0 {
			return
		}

		topCtx.InputChannels["data"] = append(topCtx.InputChannels["data"], message)
	}
}
