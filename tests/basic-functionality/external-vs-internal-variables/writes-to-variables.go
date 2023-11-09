package externalvsinternalvariablestests

import (
	"fmt"

	interpretercore "gitlab.com/jbyte777/prompt-ql/v4/core"
	interpreter "gitlab.com/jbyte777/prompt-ql/v4/interpreter"
)

// 28-09-2023: Works +++
func WriteToInternalVariableTest() {
	interpreterInst := interpreter.New(
		interpreter.TPromptQLOptions{},
	)

	result := interpreterInst.Instance.Execute(
		`
			{~set to="myVar"}Hello, PromptQL!{/set}
			My internal variable contains:
			{~get from="myVar" /}
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

// 28-09-2023: Works +++
func WriteToExternalVariableTest() {
	defaultGlobals := interpretercore.TGlobalVariablesTable{
		"myVar": "!",
	}
	interpreterInst := interpreter.New(
		interpreter.TPromptQLOptions{
			DefaultExternalGlobals: defaultGlobals,
		},
	)

	result := interpreterInst.Instance.Execute(
		`
			{~set to=@myVar}Hello, PromptQL!{/set}
			My external variable contains:
			{~get from=@myVar /}
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
