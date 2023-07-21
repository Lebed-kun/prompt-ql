package basicfunctionalitytests

import (
	"fmt"

	interpreter "gitlab.com/jbyte777/prompt-ql/interpreter"
)

// Works ++++
func InterpreterErrorTest(
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
			{~set ="X"}Example text{/set}
			{~get ="X" /}
			Hello world!
		`,
		nil,
	)

	if result.Error != nil {
		panic(result.Error)
	}

	fmt.Println("===================")
	resultStr, _ := result.ResultDataStr()
	fmt.Printf(
		"ChatGPT response:\n%v",
		resultStr,
	)
	fmt.Println("===================")
}
