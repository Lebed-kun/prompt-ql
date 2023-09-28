package basicfunctionalitytests

import (
	"fmt"

	interpreter "gitlab.com/jbyte777/prompt-ql/interpreter"
)

// Works ++++++
func ExecutionErrorTest() {
	interpreterInst := interpreter.New(
		interpreter.TPromptQLOptions{},
	)

	result := interpreterInst.Instance.Execute(
		`
			{~call fn="nonexistentfn"}Example text{/call}
		`,
	)

	if result.Error != nil {
		panic(result.Error)
	}

	fmt.Println("===================")
	resultStr, _ := result.ResultDataStr()
	errStr, _ := result.ResultErrorStr()
	fmt.Printf(
		"ChatGPT response:\n%v",
		resultStr,
	)
	fmt.Printf(
		"ChatGPT error:\n%v\n",
		errStr,
	)
	fmt.Println("===================")
}
