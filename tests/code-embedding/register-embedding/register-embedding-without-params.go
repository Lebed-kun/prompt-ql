package registerembeddingtest

import (
	"fmt"

	interpreter "gitlab.com/jbyte777/prompt-ql/v5/interpreter"
)

// 11-11-2023: Works +++
// 12-11-2023: regress +++
func RegisterEmbeddingWithoutParamsTest() {
	interpreterInst := interpreter.New(
		interpreter.TPromptQLOptions{},
	)

	resultOnDefinition := interpreterInst.Instance.Execute(
		`
			{~session_begin /}
			{~unsafe_clear_stack /}
			{~embed_def name="myEmbedding1"}
				<%
					{~set to="myVar1"}
						{~data}Hello, Alice!{/data}
					{/set}
					{~get from="myVar1" /}
				%>
			{/embed_def}
			{~data}Just an ordinary command executed on definition phase{/data}
			{~embed_def name="myEmbedding2"}
				<%
					{~set to="myVar2"}
						{~data}Hello, Bob!{/data}
					{/set}
					{~get from="myVar2" /}
				%>
			{/embed_def}
		`,
	)
	if resultOnDefinition.Error != nil {
		panic(resultOnDefinition.Error)
	}
	resultOnDefStr, _ := resultOnDefinition.ResultDataStr()

	resultOnExpansion := interpreterInst.Instance.Execute(
		`
			{~unsafe_clear_stack /}
			{~embed_exp name="myEmbedding1" /}
			{~embed_exp name="myEmbedding2" /}
		`,
	)
	if resultOnExpansion.Error != nil {
		panic(resultOnExpansion.Error)
	}
	resultOnExpStr, _ := resultOnExpansion.ResultDataStr()

	resultOnExecution := interpreterInst.Instance.Execute(
		fmt.Sprintf(`
			{~unsafe_clear_stack /}
			%v
		`,
		resultOnExpStr,
		),
	)
	if resultOnExecution.Error != nil {
		panic(resultOnExecution.Error)
	}
	resultOnExecStr, _ := resultOnExecution.ResultDataStr()

	fmt.Println("===================")
	fmt.Printf(
		"Agent response on definition phase:\n%v\n=====\n",
		resultOnDefStr,
	)
	fmt.Printf(
		"Agent response on expansion phase:\n%v\n=====\n",
		resultOnExpStr,
	)
	fmt.Printf(
		"Agent response on execution phase:\n%v\n=====\n",
		resultOnExecStr,
	)

	fmt.Println("===================")
}
