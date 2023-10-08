package externalvsinternalvariablestests

import (
	"fmt"

	interpretercore "gitlab.com/jbyte777/prompt-ql/v2/core"
	interpreter "gitlab.com/jbyte777/prompt-ql/v2/interpreter"
)

// 28-09-2023: Works +++
func InternalVariableStateAfterSessionTest() {
	interpreterInst := interpreter.New(
		interpreter.TPromptQLOptions{},
	)

	result1 := interpreterInst.Instance.Execute(
		`
			{~set to="myVar"}Hello, PromptQL!{/set}
			1st session:
			{~get from="myVar" /}
		`,
	)
	result2 := interpreterInst.Instance.Execute(
		`
			2nd session:
			{~get from="myVar" /}
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
		"GPT 1st response:\n%v\n",
		result1Str,
	)
	fmt.Printf(
		"GPT 2nd response:\n%v\n",
		result2Str,
	)
	fmt.Println("===================")
}

// 28-09-2023: Works +++
func ExternalVariableStateAfterSessionTest() {
	defaultGlobals := interpretercore.TGlobalVariablesTable{
		"myVar": "Hello, PromptQL!",
	}
	interpreterInst := interpreter.New(
		interpreter.TPromptQLOptions{
			DefaultExternalGlobals: defaultGlobals,
		},
	)

	result1 := interpreterInst.Instance.Execute(
		`
			1st session:
			{~get from=@myVar /}
		`,
	)
	result2 := interpreterInst.Instance.Execute(
		`
			2nd session:
			{~get from=@myVar /}
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
		"GPT 1st response:\n%v\n",
		result1Str,
	)
	fmt.Printf(
		"GPT 2nd response:\n%v\n",
		result2Str,
	)
	fmt.Println("===================")
}
