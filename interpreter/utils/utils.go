package interpreterutils

import (
	"fmt"
	"strings"
	api "gitlab.com/jbyte777/prompt-ql/v2/api"
)

func MergeGptApiChoices(choices []api.TGptApiResponseChoice) string {
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
