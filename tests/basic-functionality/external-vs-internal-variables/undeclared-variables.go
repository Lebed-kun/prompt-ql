package externalvsinternalvariablestests

import (
	"fmt"

	interpretercore "gitlab.com/jbyte777/prompt-ql/v4/core"
	interpreter "gitlab.com/jbyte777/prompt-ql/v4/interpreter"
)

// 28-09-2023: Works +++
func InternalUndeclaredVariableTest() {
	interpreterInst := interpreter.New(
		interpreter.TPromptQLOptions{},
	)

	result := interpreterInst.Instance.Execute(
		`
			{~set to="myVar222"}Hello, PromptQL!{/set}
			My internal undeclared variable contains:
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
func ExternalUndeclaredVariableTest() {
	defaultGlobals := interpretercore.TGlobalVariablesTable{
		"myVar222": "Hello, PromptQL!",
	}
	interpreterInst := interpreter.New(
		interpreter.TPromptQLOptions{
			DefaultExternalGlobals: defaultGlobals,
		},
	)

	result := interpreterInst.Instance.Execute(
		`
			My external undeclared variable contains:
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
