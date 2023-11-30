package inlineembeddingtest

import (
	"fmt"

	interpreter "gitlab.com/jbyte777/prompt-ql/v5/interpreter"
)

// 12-11-2023: Works +++
func ExpandInlineEmbeddingTest() {
	interpreterInst := interpreter.New(
		interpreter.TPromptQLOptions{},
	)

	// expansion
	resultOnExpansion := interpreterInst.Instance.Execute(
		`
			{~session_begin /}
			{~unsafe_clear_stack /}
			<%
				{~set to="myVar1"}Hello, Alice!{/set}
				{~get from="myVar1" /}
			%>
			{~embed_exp inline}
				<%
					{~set to="myVar2"}Hello, Bob!{/set}
					{~get from="myVar2" /}
					%someothercmd
				%>
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
		"Agent response on expansion phase:\n\"%v\"\n=====\n",
		resultOnExpStr,
	)
	fmt.Printf(
		"Agent response on execution phase:\n\"%v\"\n=====\n",
		resultOnExecStr,
	)
	fmt.Println("===================")
}
