package externalvsinternalvariablestests

import (
	"fmt"

	interpretercore "gitlab.com/jbyte777/prompt-ql/v3/core"
	interpreter "gitlab.com/jbyte777/prompt-ql/v3/interpreter"
)

// 28-09-2023: Works +++
func ExternalVariableDereferenceTest() {
	defaultGlobals := interpretercore.TGlobalVariablesTable{
		"myRef": "@myVar",
		"myVar": "Hello, PromptQL!",
	}
	interpreterInst := interpreter.New(
		interpreter.TPromptQLOptions{
			DefaultExternalGlobals: defaultGlobals,
		},
	)

	result := interpreterInst.Instance.Execute(
		`
			My external referenced variable contains:
			{~get from=$@myRef /}
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
func ExternalEmptyVariableDereferenceErrorTest() {
	defaultGlobals := interpretercore.TGlobalVariablesTable{
		"myRef": "@myVar",
		"myVar": "Hello, PromptQL!",
	}
	interpreterInst := interpreter.New(
		interpreter.TPromptQLOptions{
			DefaultExternalGlobals: defaultGlobals,
		},
	)

	result := interpreterInst.Instance.Execute(
		`
			My external referenced variable contains:
			{~get from=$@ /}
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
func ExternalVariableAtCharIsPartOfStringLiteralTest() {
	interpreterInst := interpreter.New(
		interpreter.TPromptQLOptions{},
	)

	result := interpreterInst.Instance.Execute(
		`
			{~set to="query1"}
				{~data}
					{~open_query sync model=@}
						Tell me something
					{/open_query}
				{/data}
			{/set}
			{~set to="query2"}
				{~data}
					{~open_query sync model=@myModel}
						Tell me your name
					{/open_query}
				{/data}
			{/set}

			Error 1st:
			{~get from="query1" /}
			===========
			Error 2nd:
			{~get from="query2" /}
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
func ExternalVariableSameLiteralRepresentationTest() {
	defaultGlobals := interpretercore.TGlobalVariablesTable{
		"myVar": "Hello, PromptQL!",
	}
	interpreterInst := interpreter.New(
		interpreter.TPromptQLOptions{
			DefaultExternalGlobals: defaultGlobals,
		},
	)

	result := interpreterInst.Instance.Execute(
		`
			My external variable contains this:
			{~get from=@myVar /}
			==================================
			And it also contains this:
			{~get from="@myVar" /}
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
