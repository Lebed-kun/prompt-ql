package basicfunctionalitytests

import (
	"fmt"

	interpreter "gitlab.com/jbyte777/prompt-ql/v5/interpreter"
)

// Works +++
// 28-09-2023: Works on total regress +++
// 11-11-2023: Works on regress +++
// 30-11-2023: random regress +++
func NestedWrappersTest() {
	interpreterInst := interpreter.New(
		interpreter.TPromptQLOptions{},
	)

	result := interpreterInst.Instance.Execute(
		`
			{~data}
				{~user}Example text{/user}
			{/data}
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
