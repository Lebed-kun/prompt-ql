package interpreterutils

import (
	b64 "encoding/base64"
	"errors"
	"fmt"
	"strings"

	chatapi "gitlab.com/jbyte777/prompt-ql/v5/default-apis/chat-api"
	ttiapi "gitlab.com/jbyte777/prompt-ql/v5/default-apis/tti-api"
	httputils "gitlab.com/jbyte777/prompt-ql/v5/utils/http"
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


func ReadTtiBlob(
	ttiResponse *ttiapi.TTtiApiResponse,
) ([]byte, error) {
	if len(ttiResponse.Data) < 1 {
		return nil, errors.New("tti response does not contain image")
	}

	imgEntryUrl := ttiResponse.Data[0].Url
	imgEntryB64 := ttiResponse.Data[0].B64Json

	if len(imgEntryUrl) > 0 {
		imgData, err := httputils.DoHttpRequest(
			imgEntryUrl,
			"GET",
			"",
			30,
		)
		if err != nil {
			return nil, err
		}

		return imgData, nil
	} else if len(imgEntryB64) > 0 {
		// TODO: investigate b64_json format
		blob, err := b64.StdEncoding.DecodeString(imgEntryB64)
		if err != nil {
			return nil, err
		}

		return blob, nil
	}

	return nil, errors.New("no img data returned")
}
