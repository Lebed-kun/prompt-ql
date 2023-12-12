package basicfunctionalitytests

import (
	"fmt"

	interpreter "gitlab.com/jbyte777/prompt-ql/v5/interpreter"
)

// Works +++
// 28-09-2023: Works on total regress +++
func MismatchCommandErrorTest() {
	interpreterInst := interpreter.New(
		interpreter.PromptQLOptions{},
	)

	result := interpreterInst.Instance.Execute(
		`
			{~data}
				{~user}Example text{/user}
			{/call}
		`,
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
