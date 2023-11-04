package hellocommandtests

import (
	"fmt"
	interpreter "gitlab.com/jbyte777/prompt-ql/v3/interpreter"
)

// 04-11-2023: Works +++
func HeaderCommandTest() {
	interpreterInst := interpreter.New(
		interpreter.TPromptQLOptions{},
	)

	result := interpreterInst.Instance.Execute(
		`
			Message header:
			{~header from="Alice" to="Bob" /}
		`,
	)

	if result.Error != nil {
		panic(result.Error)
	}

	fmt.Println("===================")
	resultStr, _ := result.ResultDataStr()
	errStr, _ := result.ResultErrorStr()
	fmt.Printf(
		"GPT response:\n%v\n",
		resultStr,
	)
	fmt.Printf(
		"GPT error:\n%v\n",
		errStr,
	)
	fmt.Println("===================")
}
