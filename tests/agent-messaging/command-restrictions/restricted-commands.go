package commandrestrictionstests

import (
	"fmt"

	interpretercore "gitlab.com/jbyte777/prompt-ql/v5/core"
	interpreter "gitlab.com/jbyte777/prompt-ql/v5/interpreter"
)

// 11-11-2023: Works +++
// 10-12-2023: Regress ???
func RestrictedCommandsTest() {
	interpreterInst := interpreter.New(
		interpreter.TPromptQLOptions{
			RestrictedCommands: interpretercore.TRestrictedCommands{
				"unsafe_clear_stack": true,
				"embed_def": true,
			},
		},
	)

	result1 := interpreterInst.Instance.Execute(
		`
			{~hello /}
	    Then I try to execute forbidden command for clearing stack:
			{~unsafe_clear_stack /}
		`,
	)
	result1Str, _ := result1.ResultDataStr()

	result2 := interpreterInst.Instance.Execute(
		`
			{~session_begin /}
			{~hello /}
	    Then I try to execute forbidden command for embedding definition:
			{~embed_def name="myEmbedding"}<% {~hello /} %>{/embed_def}
		`,
	)
	result2Str, _ := result2.ResultDataStr()

	interpreterInst.Instance.UnsafeExecute(`{~unsafe_clear_stack /}`)

	result3 := interpreterInst.Instance.Execute(
		`
			I try to expand defined embedding:
			{~embed_exp name="myEmbedding" /}
		`,
	)
	result3Str, _ := result3.ResultDataStr()

	fmt.Println("===================")
	fmt.Printf(
		"Agent result one:\n%v\n",
		result1Str,
	)
	fmt.Printf(
		"Agent result two:\n%v\n",
		result2Str,
	)
	fmt.Printf(
		"Agent result three:\n%v\n",
		result3Str,
	)
	if result1.Error != nil {
		fmt.Printf(
			"Agent error one:\n%v\n",
			result1.Error.Error(),
		)
	}
	if result2.Error != nil {
		fmt.Printf(
			"Agent error two:\n%v\n",
			result2.Error.Error(),
		)
	}
	fmt.Println("===================")
}

