package basicfunctionalitytests

import (
	"fmt"

	interpreter "gitlab.com/jbyte777/prompt-ql/v5/interpreter"
)

// Works ++++++
// 28-09-2023: Works on total regress +++
func ExecutionErrorTest() {
	interpreterInst := interpreter.New(
		interpreter.PromptQLOptions{},
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
