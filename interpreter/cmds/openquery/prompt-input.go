package openquerycmd

import (
	"fmt"

	api "gitlab.com/jbyte777/prompt-ql/api"
	interpreter "gitlab.com/jbyte777/prompt-ql/core"
	parseutils "gitlab.com/jbyte777/prompt-ql/utils/parse"
	cmdtypes "gitlab.com/jbyte777/prompt-ql/interpreter/cmds/types"
	mathutils "gitlab.com/jbyte777/prompt-ql/utils/math"
)

const maxMsgId int = 1_000

func getPromptInfo(
	inputs interpreter.TFunctionInputChannelTable,
	msgRole string,
	ptr int,
) *cmdtypes.TPromptInfo {
	if msgChan, hasMsgChan := inputs[msgRole]; hasMsgChan {
		if ptr < len(msgChan) {
			rawMsg := msgChan[ptr].(string)
			msgIdx, err, msgBegin := parseutils.ParseUintFromPrefix(rawMsg)
			if err != nil {
				return &cmdtypes.TPromptInfo{
					MsgId: maxMsgId,
					Err: fmt.Errorf(
						"!error %v",
						err.Error(),
					),
					MsgBegin: -1,
				}
			}

			return &cmdtypes.TPromptInfo{
				MsgId: msgIdx,
				Err: nil,
				MsgBegin: msgBegin,
			}
		}
	}

	return &cmdtypes.TPromptInfo{
		MsgId: maxMsgId,
		Err: nil,
		MsgBegin: -1,
	}
}

func getPrompts(
	inputs interpreter.TFunctionInputChannelTable,
) ([]api.TMessage, error) {
	res := make([]api.TMessage, 0)

	userPtr := 0
	assistantPtr := 0
	systemPtr := 0
	for {
		userPromptInfo := getPromptInfo(inputs, "user", userPtr)
		if userPromptInfo.Err != nil {
			return nil, userPromptInfo.Err
		}

		assistantPromptInfo := getPromptInfo(inputs, "assistant", assistantPtr)
		if assistantPromptInfo.Err != nil {
			return nil, assistantPromptInfo.Err
		}

		systemPromptInfo := getPromptInfo(inputs, "system", systemPtr)
		if systemPromptInfo.Err != nil {
			return nil, systemPromptInfo.Err
		}

		minId := mathutils.MinInt(
			userPromptInfo.MsgId,
			mathutils.MinInt(
				assistantPromptInfo.MsgId,
				systemPromptInfo.MsgId,
			),
		)

		if minId == maxMsgId {
			break
		}

		var msg api.TMessage
		switch (minId) {
		case userPromptInfo.MsgId:
			rawMsg := inputs["user"][userPtr].(string) 
			msgContent := rawMsg[userPromptInfo.MsgBegin + 1:]

			msg = api.TMessage{
				Role: "user",
				Content: msgContent,
			}
			userPtr++
		case assistantPromptInfo.MsgId:
			rawMsg := inputs["assistant"][assistantPtr].(string) 
			msgContent := rawMsg[assistantPromptInfo.MsgBegin + 1:]

			msg = api.TMessage{
				Role: "assistant",
				Content: msgContent,
			}
			assistantPtr++
		case systemPromptInfo.MsgId:
			rawMsg := inputs["system"][systemPtr].(string) 
			msgContent := rawMsg[systemPromptInfo.MsgBegin + 1:]

			msg = api.TMessage{
				Role: "system",
				Content: msgContent,
			}
			systemPtr++
		}

		res = append(res, msg)
	}

	if len(res) == 0 {
		return nil, fmt.Errorf(
			"!error prompts are empty",
		)
	}

	return res, nil
}
