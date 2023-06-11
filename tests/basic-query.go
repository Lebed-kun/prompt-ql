package tests

import (
	"fmt"

	interpreter "gitlab.com/jbyte777/prompt-ql/interpreter"
)

func BasicQueryTest(
	openAiBaseUrl string,
	openAiKey string,
) {
	interpreterInst := interpreter.NewInterpreter(
		openAiBaseUrl,
		openAiKey,
	)

	result := interpreterInst.ExecuteClean(
		`
			{~open_query to="query1" model="gpt-4"}
				{~system}
					You are a helpful and terse assistant.
				{/system}
				I want a response to the following question:
				Name 3 world-class experts (past or present) who would be great at answering this?
				Don't answer the question yet.
			{/open_query}
			{~listen_query from="query1" /}
		`,
	)

	if result.Error != nil {
		panic(result.Error)
	}

	fmt.Println("===================")
	fmt.Printf(
		"ChatGPT response:\n%v",
		result.Result["data"],
	)
	fmt.Println("===================")
}
