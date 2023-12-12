package basicfunctionalitytests

import (
	"fmt"

	interpretercore "gitlab.com/jbyte777/prompt-ql/v5/core"
	interpreter "gitlab.com/jbyte777/prompt-ql/v5/interpreter"
)

// Works ++++
// 28-09-2023: Works on total regress +++
// 04-11-2023: random regression test +++
// 11-11-2023: random regression test +++
// 30-11-2023: random regression test +++
func QueryWithPostprocessFunctionTest(
	openAiBaseUrl string,
	openAiKey string,
) {
	defaultGlobals := interpretercore.TGlobalVariablesTable{
		"postprocess": postProcessFunctionTest,
	}
	interpreterInst := interpreter.New(
		interpreter.PromptQLOptions{
			OpenAiBaseUrl: openAiBaseUrl,
			OpenAiKey: openAiKey,
			DefaultExternalGlobals: defaultGlobals,
		},
	)

	result := interpreterInst.Instance.Execute(
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
			{~call fn=@postprocess }
				{~get from="queryres" /}
			{/call}
		`,
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
