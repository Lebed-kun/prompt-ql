package registerembeddingtest

import (
	"fmt"

	interpretercore "gitlab.com/jbyte777/prompt-ql/v5/core"
	interpreter "gitlab.com/jbyte777/prompt-ql/v5/interpreter"
)

// 11-11-2023: Works +++
func AgentLayoutWithEmbeddingTest() {
	// Define external globals
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

	// Define custom ML APIs
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

	// Define custom code embeddings
	interpreterInst.Instance.Execute(
		`
			{~session_begin /}
			{~embed_def name="myEmbedding1" desc="Alice code embedding"}
				<%
					{~set to="myVar1"}
						{~data}Hello, Alice!{/data}
					{/set}
					{~get from="myVar1" /}
					%someothercmd
				%>
			{/embed_def}
		`,
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
