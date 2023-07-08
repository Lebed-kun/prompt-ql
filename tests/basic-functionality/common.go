package basicfunctionalitytests

import (
	"encoding/json"
	"fmt"
	"strings"
)

type tPostProcessStepTest struct {
	Num int `json:"num"`
	Text string `json:"text"`
}

type tPostProcessResultTest struct {
	Steps []tPostProcessStepTest `json:"steps"`
}

func isWhitespace(r rune) bool {
	return r == '\n' ||
		r == '\r' ||
		r == '\t' ||
		r == ' '
}

func postProcessFunctionTest(args []interface{}) interface{} {
	if len(args) < 1 {
		return "!error arguments list for postprocessing is empty"
	}

	text, isTextStr := args[0].(string)
	if !isTextStr {
		return "!error text for postprocessing is not string"
	}

	steps := make([]tPostProcessStepTest, 0)
	currStep := 0
	currText := strings.Builder{}
	textRunes := []rune(text)
	textPtr := 0
	for textPtr < len(textRunes) {
		if isWhitespace(textRunes[textPtr]) {
			for textPtr < len(textRunes) && isWhitespace(textRunes[textPtr]) {
				textPtr++
			}
		}

		if textPtr < len(textRunes) && '0' <= textRunes[textPtr] && textRunes[textPtr] <= '9' {
			num := 0
			for textPtr < len(textRunes) && '0' <= textRunes[textPtr] && textRunes[textPtr] <= '9' {
				dig := int(textRunes[textPtr] - '0')
				num = num * 10 + dig
				textPtr++
			}

			if textRunes[textPtr] == '.' {
				textPtr++
			}

			currStep = num
			continue
		}

		for textPtr < len(textRunes) && textRunes[textPtr] != '\n' {
			currText.WriteRune(textRunes[textPtr])
			textPtr++
		}
		textPtr++
		
		cleanText := strings.TrimSpace(currText.String())
		if currStep != 0 && len(cleanText) > 0 {
			steps = append(steps, tPostProcessStepTest{
				Num: currStep,
				Text: cleanText,
			})
			currStep = 0
			currText = strings.Builder{}
		}
	}

	postProcessResult := tPostProcessResultTest{
		Steps: steps,
	}
	postProcessResultStr, _ := json.MarshalIndent(
		postProcessResult,
		" ",
		" ",
	)

	return fmt.Sprintf("!data %v", string(postProcessResultStr))
}