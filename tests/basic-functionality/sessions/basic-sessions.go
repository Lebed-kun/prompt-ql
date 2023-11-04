package sessionstests

import (
	"fmt"
	interpreter "gitlab.com/jbyte777/prompt-ql/v3/interpreter"
)

// 07-10-2023: Works +++
func BasicSession() {
	interpreterInst := interpreter.New(
		interpreter.TPromptQLOptions{},
	)

	interpreterInst.Instance.Execute(
		`
			{~session_begin /}
			{~set to="myVar"}Hello, PromptQL!{/set}
		`,
	)
	result1 := interpreterInst.Instance.Execute(
		`
	    Value of "myVar" during session is:
			{~get from="myVar" /}
			==================
			{~session_end /}
		`,
	)
	result2 := interpreterInst.Instance.Execute(
		`
	    Value of "myVar" after session is:
			{~get from="myVar" /}
			=================
		`,
	)

	if result1.Error != nil {
		panic(result1.Error)
	}
	if result2.Error != nil {
		panic(result2.Error)
	}

	fmt.Println("===================")
	result1Str, _ := result1.ResultDataStr()
	result2Str, _ := result2.ResultDataStr()
	fmt.Printf(
		"GPT first result:\n%v\n",
		result1Str,
	)
	fmt.Printf(
		"GPT second result:\n%v\n",
		result2Str,
	)
	fmt.Println("===================")
}
