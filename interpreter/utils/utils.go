package interpreterutils

import (
	"fmt"
	"strings"
	chatapi "gitlab.com/jbyte777/prompt-ql/v5/default-apis/chat-api"
)

func MergeGptApiChoices(choices []chatapi.TGptApiResponseChoice) string {
	strChoices := make([]string, 0)
	for _, choice := range choices {
		if choice.Message.Role != "assistant" {
			continue
		}

		strChoices = append(strChoices, choice.Message.Content)
	}

	result := fmt.Sprintf(
		"!assistant %v",
		strings.Join(strChoices, "\n=====\n"),
	)
	return result
}
