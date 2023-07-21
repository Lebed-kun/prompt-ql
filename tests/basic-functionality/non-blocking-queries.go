package basicfunctionalitytests

import (
	"fmt"

	interpretercore "gitlab.com/jbyte777/prompt-ql/core"
	interpreter "gitlab.com/jbyte777/prompt-ql/interpreter"
	timeutils "gitlab.com/jbyte777/prompt-ql/utils/time"
)

func logTimeForProgram(args []interface{}) interface{} {
	if len(args) < 1 {
		return ""
	}

	log, isLogStr := args[0].(string)
	if !isLogStr {
		return ""
	}

	fmt.Printf(
		"[%v] %v",
		timeutils.NowTimestamp(),
		log,
	)

	return ""
}

// Works ++++
func NonBlockingQueriesTest(
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
				Write a comprehensive guide to learn statistics step by step.
			{/open_query}
			{~call fn="logtime"}
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
			{~call fn="logtime"}
				open query2
			{/call}
			=======================
			Answer1: {~listen_query from="query1" /}
			{~call fn="logtime"}
				listen query1
			{/call}
			=======================
			Answer2: {~listen_query from="query2" /}
			{~call fn="logtime"}
				listen query2
			{/call}
			=======================
		`,
		interpretercore.TGlobalVariablesTable{
			"logtime": logTimeForProgram,
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
