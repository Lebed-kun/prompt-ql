package wrapperswicth

import (
	interpreter "gitlab.com/jbyte777/prompt-ql/core"
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
		topCtx.InputChannels["data"] = append(
			topCtx.InputChannels["data"],
			data,
		)
	}
}
