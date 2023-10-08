package hellocommandtests

import (
	"fmt"

	interpretercore "gitlab.com/jbyte777/prompt-ql/v2/core"
	interpreter "gitlab.com/jbyte777/prompt-ql/v2/interpreter"
)

// 07-10-2023: Works +++
func HelloCommandTest() {
	defaultGlobals := interpretercore.TGlobalVariablesTable{
		"myRef": "@myVar",
		"myVar": "Hello, PromptQL!",
		"myFunc": func(args []interface{}) interface{} {
			return nil
		},
	}
	interpreterInst := interpreter.New(
		interpreter.TPromptQLOptions{
			DefaultExternalGlobals: defaultGlobals,
		},
	)

	interpreterInst.CustomApis.RegisterModelApi(
		"myMlModel111",
		func(
			model string,
			temperature float64,
			inputs interpretercore.TFunctionInputChannelTable,
			execInfo interpretercore.TExecutionInfo,
		) (string, error) {
			return "", nil
		},
	)
	interpreterInst.CustomApis.RegisterModelApi(
		"myMlModel222",
		func(
			model string,
			temperature float64,
			inputs interpretercore.TFunctionInputChannelTable,
			execInfo interpretercore.TExecutionInfo,
		) (string, error) {
			return "", nil
		},
	)

	result := interpreterInst.Instance.Execute(
		`
			My agent layout is:
			{~hello /}
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