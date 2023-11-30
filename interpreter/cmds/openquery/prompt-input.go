package openquerycmd

import (
	"fmt"

	api "gitlab.com/jbyte777/prompt-ql/v5/api"
	interpreter "gitlab.com/jbyte777/prompt-ql/v5/core"
)

func getPrompts(
	inputs interpreter.TFunctionInputChannelTable,
	execInfo interpreter.TExecutionInfo,
) ([]api.TMessage, error) {
	res := make([]api.TMessage, 0)

	userChan := inputs["user"]
	assistantChan := inputs["assistant"]
	systemChan := inputs["system"]

	ptr := 0
	for {
		var msg api.TMessage
		userChanMsg := ""
		assistantChanMsg := ""
		systemChanMsg := ""

		if ptr >= len(userChan) && ptr >= len(assistantChan) && ptr >= len(systemChan) {
			break
		}
		
		if ptr < len(userChan)  {
			userChanMsg = userChan[ptr].(string)
		}
		if ptr < len(assistantChan) {
			assistantChanMsg = assistantChan[ptr].(string)
		}
		if ptr < len(systemChan) {
			systemChanMsg = systemChan[ptr].(string)
		}

		if len(userChanMsg) > 0 {
			msg = api.TMessage{
				Role: "user",
				Content: userChanMsg,
			}
		} else if len(assistantChanMsg) > 0 {
			msg = api.TMessage{
				Role: "assistant",
				Content: assistantChanMsg,
			}
		} else if len(systemChanMsg) > 0 {
			msg = api.TMessage{
				Role: "system",
				Content: systemChanMsg,
			}
		}

		res = append(res, msg)
		ptr++
	}

	if len(res) == 0 {
		return nil, fmt.Errorf(
			"!error (line=%v, char=%v): prompts are empty",
			execInfo.Line,
			execInfo.CharPos,
		)
	}

	return res, nil
}
