package tests

import (
	"encoding/json"
	"fmt"
	"strings"

	interpretercore "gitlab.com/jbyte777/prompt-ql/core"
	interpreter "gitlab.com/jbyte777/prompt-ql/interpreter"
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

func QueryWithPostprocessFunctionTest(
	openAiBaseUrl string,
	openAiKey string,
) {
	interpreterInst := interpreter.New(
		openAiBaseUrl,
		openAiKey,
	)

	result := interpreterInst.Execute(
		`
			{~open_query to="query1" model="gpt-3.5-turbo-16k"}
				{~system}
					You are a helpful and terse assistant.
				{/system}
				I want a response to the following question:
				Write a comprehensive guide to machine learning step by step
			{/open_query}
			{~set to="queryres"}
				{~listen_query from="query1" /}
			{/set}
			Raw result is:
			{~get from="queryres" /}

			JSON result is:
			{~call fn="postprocess"}
				{~get from="queryres" /}
			{/call}
		`,
		interpretercore.TGlobalVariablesTable{
			"postprocess": postProcessFunctionTest,
		},
	)
	/*
	result := interpreterInst.Execute(
		`
			{~set to="queryres"}
				1. Make a dish
				2. Eat it
				3. ...
			{/set}
			Raw result is:
			{~get from="queryres" /}
			==========++++++==========
			JSON result is:
			{~call fn="postprocess"}
				{~get from="queryres" /}
			{/call}
		`,
		interpretercore.TGlobalVariablesTable{
			"postprocess": postProcessFunctionTest,
		},
	)
	*/
	/*
	result := interpreterInst.Execute(
		`
			{~call fn="postprocess"}
				1. Make a dish
				2. Eat it
				3. ...
			{/call}
		`,
		interpretercore.TGlobalVariablesTable{
			"postprocess": postProcessFunctionTest,
		},
	)
	*/

	if result.Error != nil {
		panic(result.Error)
	}

	fmt.Println("===================")
	resultStr, _ := result.ResultDataStr()
	errStr, _ := result.ResultErrorStr()
	fmt.Printf(
		"ChatGPT response:\n%v\n",
		resultStr,
	)
	fmt.Printf(
		"ChatGPT error:\n%v\n",
		errStr,
	)
	fmt.Println("===================")
}
