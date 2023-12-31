package syncqueriestests

import (
	"fmt"

	interpreter "gitlab.com/jbyte777/prompt-ql/interpreter"
	testutils "gitlab.com/jbyte777/prompt-ql/tests/utils"
	interpretercore "gitlab.com/jbyte777/prompt-ql/core"
)

// Works +
func BasicSyncQueryTest(
	openAiBaseUrl string,
	openAiKey string,
) {
	interpreterInst := interpreter.New(
		openAiBaseUrl,
		openAiKey,
		0,
		0,
	)

	result := interpreterInst.Instance.Execute(
		`
			=========^^^^^==========
			{~call fn="logtime"}
				begin first query
			{/call}
			{~open_query sync model="gpt-3.5-turbo-16k"}
				{~system}
					You are a helpful and terse assistant.
				{/system}
				I want a response to the following question:
				Write a comprehensive guide to write an article for Medium
			{/open_query}
			{~call fn="logtime"}
				end first query
			{/call}
			=========+++++==========
			{~call fn="logtime"}
				begin second query
			{/call}
			{~open_query sync model="gpt-3.5-turbo-16k"}
				{~system}
					You are a helpful and terse assistant.
				{/system}
				I want a response to the following question:
				Write a comprehensive guide to write an order email
			{/open_query}
			{~call fn="logtime"}
				end second query
			{/call}
			=========^^^^^==========
		`,
		interpretercore.TGlobalVariablesTable{
			"logtime": testutils.LogTimeForProgram,
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
