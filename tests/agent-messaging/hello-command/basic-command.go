package hellocommandtests

import (
	"fmt"

	interpretercore "gitlab.com/jbyte777/prompt-ql/v5/core"
	interpreter "gitlab.com/jbyte777/prompt-ql/v5/interpreter"
)

// 07-10-2023: Works +++
// 04-11-2023: Works with descriptions +++
// 11-11-2023: Works on regress +++
// 23-12-2023: Works on regress +++
func HelloCommandTest() {
	defaultGlobals := interpretercore.TGlobalVariablesTable{
		"myRef": "@myVar",
		"myVar": "Hello, PromptQL!",
		"myFunc": func(args []interface{}) interface{} {
			return nil
		},
	}
	defaultGlobalsMeta := interpretercore.TExternalGlobalsMetaTable{
		"myRef": &interpretercore.TExternalGlobalMetaInfo{
			Description: "a reference to @myVar",
		},
		"myVar": &interpretercore.TExternalGlobalMetaInfo{
			Description: "an example external variable",
		},
		"myFunc": &interpretercore.TExternalGlobalMetaInfo{
			Description: "an example external function",
		},
	}
	interpreterInst := interpreter.New(
		interpreter.PromptQLOptions{
			DefaultExternalGlobals: defaultGlobals,
			DefaultExternalGlobalsMeta: defaultGlobalsMeta,
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
		"my model 111",
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
		"my model 222",
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
