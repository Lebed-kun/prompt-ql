package sessionstests

import (
	"fmt"
	interpreter "gitlab.com/jbyte777/prompt-ql/v4/interpreter"
)

// 07-10-2023: Works +++
func OpenSessionErrorTest() {
	interpreterInst := interpreter.New(
		interpreter.TPromptQLOptions{},
	)

	result := interpreterInst.Instance.Execute(
		`
			{~session_begin /}
			{~session_begin /}
		`,
	)

	if result.Error != nil {
		panic(result.Error)
	}

	fmt.Println("===================")
	errorStr, _ := result.ResultErrorStr()
	fmt.Printf(
		"GPT error:\n%v\n",
		errorStr,
	)
	fmt.Println("===================")
}

// 07-10-2023: Works +++
func CloseSessionErrorTest() {
	interpreterInst := interpreter.New(
		interpreter.TPromptQLOptions{},
	)

	result := interpreterInst.Instance.Execute(
		`
			{~session_begin /}
			{~session_end /}
			{~session_end /}
		`,
	)

	if result.Error != nil {
		panic(result.Error)
	}

	fmt.Println("===================")
	errorStr, _ := result.ResultErrorStr()
	fmt.Printf(
		"GPT error:\n%v\n",
		errorStr,
	)
	fmt.Println("===================")
}
