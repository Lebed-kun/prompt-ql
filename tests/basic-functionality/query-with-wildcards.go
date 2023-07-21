package basicfunctionalitytests

import (
	"fmt"

	interpretercore "gitlab.com/jbyte777/prompt-ql/core"
	interpreter "gitlab.com/jbyte777/prompt-ql/interpreter"
)

// Works ++++
func QueryWithWildcardsTest(
	openAiBaseUrl string,
	openAiKey string,
) {
	interpreterInst := interpreter.New(
		openAiBaseUrl,
		openAiKey,
		0,
	)

	result := interpreterInst.Instance.Execute(
		`
			{~open_query to="query1" model="gpt-3.5-turbo-16k"}
				{~system}
					You are a helpful and terse assistant.
				{/system}
				I want a response to the following question:
				Write a comprehensive guide to machine learning
			{/open_query}
			{~$cmd $cmdarg=$cmdval /}
		`,
		interpretercore.TGlobalVariablesTable{
			"cmd": "listen_query",
			"cmdarg": "from",
			"cmdval": "query1",
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
