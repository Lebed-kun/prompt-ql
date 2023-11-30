package customllmstests

import (
	"fmt"
	"os/exec"
	"strings"

	interpretercore "gitlab.com/jbyte777/prompt-ql/v5/core"
	customapis "gitlab.com/jbyte777/prompt-ql/v5/custom-apis"
)

func llamaComposePrompt(
	inputs interpretercore.TFunctionInputChannelTable,
) string {
	res := make([]string, 0)

	userChan := inputs["user"]
	assistantChan := inputs["assistant"]
	systemChan := inputs["system"]

	ptr := 0
	for {
		var msg string
		userChanMsg := ""
		assistantChanMsg := ""
		systemChanMsg := ""

		if ptr >= len(userChan) && ptr >= len(assistantChan) && ptr >= len(systemChan) {
			break
		}

		if ptr < len(userChan) {
			userChanMsg = userChan[ptr].(string)
		}
		if ptr < len(assistantChan) {
			assistantChanMsg = assistantChan[ptr].(string)
		}
		if ptr < len(systemChan) {
			systemChanMsg = systemChan[ptr].(string)
		}

		if len(userChanMsg) > 0 {
			msg = userChanMsg
		} else if len(assistantChanMsg) > 0 {
			msg = assistantChanMsg
		} else if len(systemChanMsg) > 0 {
			msg = systemChanMsg
		}

		res = append(res, msg)
		ptr++
	}

	return strings.Join(res, "\n")
}

func makeLlamaDoQuery(
	pathToLlamaCommand string,
	pathToLlamaModel string,
) customapis.TDoQueryFunc {
	return func(
		model string,
		temperature float64,
		inputs interpretercore.TFunctionInputChannelTable,
		execInfo interpretercore.TExecutionInfo,
	) (string, error) {
		prompt := llamaComposePrompt(inputs)

		cmd := exec.Command(
			pathToLlamaCommand,
			"-m",
			pathToLlamaModel,
			"--temp",
			fmt.Sprintf("%.1f", temperature),
			"-p",
			fmt.Sprintf("\"%v\"", prompt),
		)

		res, err := cmd.Output()
		if err != nil {
			return "", fmt.Errorf(
				"ERROR (line=%v, charpos=%v): %v",
				execInfo.Line,
				execInfo.CharPos,
				err.Error(),
			)
		}

		return string(res), nil
	}
}
