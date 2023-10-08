package basicfunctionalitytests

import (
	"fmt"
	interpreter "gitlab.com/jbyte777/prompt-ql/v2/interpreter"
)

// Works ++++++
// 28-09-2023: Works on total regress +++
func QueryWithWildcardsTest(
	openAiBaseUrl string,
	openAiKey string,
) {
	interpreterInst := interpreter.New(
		interpreter.TPromptQLOptions{
			OpenAiBaseUrl: openAiBaseUrl,
			OpenAiKey: openAiKey,
		},
	)

	result := interpreterInst.Instance.Execute(
		`
			{~set to="cmd"}listen_query{/set}
			{~set to="cmdarg"}from{/set}
			{~set to="queryvar"}query1{/set}

			{~open_query to=$queryvar model="gpt-3.5-turbo-16k"}
				{~system}
					You are a helpful and terse assistant.
				{/system}
				I want a response to the following question:
				Write a comprehensive guide to machine learning
			{/open_query}
			{~$cmd $cmdarg=$queryvar /}
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