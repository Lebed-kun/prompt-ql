package registerembeddingtest

import (
	"fmt"

	interpreter "gitlab.com/jbyte777/prompt-ql/v5/interpreter"
)

// 11-11-2023: Works +++
// 12-11-2023: regress +++
// 30-11-2023: random regression test +++
func RegisterEmbeddingTest() {
	interpreterInst := interpreter.New(
		interpreter.PromptQLOptions{},
	)

	// definitions
	resultOnDefinition := interpreterInst.Instance.Execute(
		`
			{~session_begin /}
			{~unsafe_clear_stack /}
			{~embed_def name="myEmbedding1"}
				<%
					{~set to="myVar1"}Hello, Alice!{/set}
					{~get from="myVar1" /}
					%someothercmd
				%>
			{/embed_def}
			{~data}Just an ordinary command executed on definition phase{/data}
			{~embed_def name="myEmbedding2"}
				<%
					{~set to="myVar2"}Hello, Bob!{/set}
					{~get from="myVar2" /}
					%someothercmd
				%>
			{/embed_def}
		`,
	)
	if resultOnDefinition.Error != nil {
		panic(resultOnDefinition.Error)
	}
	resultOnDefStr, _ := resultOnDefinition.ResultDataStr()

	// expansion
	resultOnExpansion := interpreterInst.Instance.Execute(
		`
			{~unsafe_clear_stack /}
			{~embed_exp name="myEmbedding1"}
				{~data}someothercmd=<% Bob's var is: {~get from="myVar2" /} %>{/data}
			{/embed_exp}
			{~embed_exp name="myEmbedding2"}
				{~data}someothercmd=<% Alice's var is: {~get from="myVar1" /} %>{/data}
			{/embed_exp}
		`,
	)
	if resultOnExpansion.Error != nil {
		panic(resultOnExpansion.Error)
	}
	resultOnExpStr, _ := resultOnExpansion.ResultDataStr()

	// execution
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
		"Agent response on definition phase:\n\"%v\"\n=====\n",
		resultOnDefStr,
	)
	fmt.Printf(
		"Agent response on expansion phase:\n\"%v\"\n=====\n",
		resultOnExpStr,
	)
	fmt.Printf(
		"Agent response on execution phase:\n\"%v\"\n=====\n",
		resultOnExecStr,
	)
	fmt.Println("===================")
}
