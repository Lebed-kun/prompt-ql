package basicfunctionalitytests

import (
	"fmt"
	"time"

	interpretercore "gitlab.com/jbyte777/prompt-ql/core"
	interpreter "gitlab.com/jbyte777/prompt-ql/interpreter"
)

// Works ++++
func PartialExecutionTest(
	openAiBaseUrl string,
	openAiKey string,
) {
	interpreterInst := interpreter.New(
		openAiBaseUrl,
		openAiKey,
		0,
		0,
	)

	result := interpreterInst.Instance.ExecutePartial(
		`
			{~open_query to="query1" model="gpt-3.5-turbo-16k"}
				{~system}
					You are a helpful and terse assistant.
				{/system}
				I want a response to the following question:
				Write a comprehensive guide to machine learning step by step
		`,
		interpretercore.TGlobalVariablesTable{
			"postprocess": postProcessFunctionTest,
		},
	)

	time.Sleep(3 * time.Second)

	result = interpreterInst.Instance.ExecutePartial(
		`
			{/open_query}
			{~set to="queryres"}
				{~listen_query from="query1" /}
			{/set}
			Raw result is:
			{~get from="queryres" /}
			==========================
			JSON result is:
			{~call fn="postprocess"}
				{~get from="queryres" /}
			{/call}
		`,
		interpretercore.TGlobalVariablesTable{
			"postprocess": postProcessFunctionTest,
		},
	)

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
