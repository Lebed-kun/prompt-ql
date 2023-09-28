package basicfunctionalitytests

import (
	"fmt"

	interpretercore "gitlab.com/jbyte777/prompt-ql/core"
	interpreter "gitlab.com/jbyte777/prompt-ql/interpreter"
	testutils "gitlab.com/jbyte777/prompt-ql/tests/utils"
)

// Works +++++
// 28-09-2023: Works on total regress +++
func NonBlockingQueriesTest(
	openAiBaseUrl string,
	openAiKey string,
) {
	defaultGlobals := interpretercore.TGlobalVariablesTable{
		"logtime": testutils.LogTimeForProgram,
	}
	interpreterInst := interpreter.New(
		interpreter.TPromptQLOptions{
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
				Write a comprehensive guide to learn statistics step by step.
			{/open_query}
			{~call fn=@logtime }
				open query1
			{/call}
			=======================
			{~open_query to="query2" model="gpt-3.5-turbo-16k"}
				{~system}
					You are a helpful and terse assistant.
				{/system}
				I want a response to the following question:
				Write a comprehensive guide to make a solar panel step by step.
			{/open_query}
			{~call fn=@logtime }
				open query2
			{/call}
			=======================
			Answer1: {~listen_query from="query1" /}
			{~call fn=@logtime }
				listen query1
			{/call}
			=======================
			Answer2: {~listen_query from="query2" /}
			{~call fn=@logtime }
				listen query2
			{/call}
			=======================
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
