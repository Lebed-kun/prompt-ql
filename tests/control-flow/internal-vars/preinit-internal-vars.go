package internalvarstests

import (
	"fmt"

	interpretercore "gitlab.com/jbyte777/prompt-ql/v5/core"
	interpreter "gitlab.com/jbyte777/prompt-ql/v5/interpreter"
)

// 12-12-2023 Works +++
func PreinitInternalVarsTest() {
	defaultInternals := interpretercore.TGlobalVariablesTable{
		"_internal_fn": func(args []interface{}) interface{} {
			return "Hello Alice!"
		},
	}
	interpreterInst := interpreter.New(
		interpreter.PromptQLOptions{
			PreinitializedInternalGlobals: defaultInternals,
		},
	)

	result1 := interpreterInst.Instance.Execute(
		`
			First call of internal function should resolve into error:
			{~call fn="_internal_fn" /}
		`,
	)
	result2 := interpreterInst.Instance.UnsafeExecute(
		`
			{~unsafe_preinit_vars /}
			Second call of internal function should be ok:
			{~call fn="_internal_fn" /}
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
	err1Str, _ := result1.ResultErrorStr()
	fmt.Printf(
		"Agent response 1:\n%v\n",
		result1Str,
	)
	fmt.Printf(
		"Agent error 1:\n%v\n",
		err1Str,
	)

	result2Str, _ := result2.ResultDataStr()
	err2Str, _ := result2.ResultErrorStr()
	fmt.Printf(
		"Agent response 2:\n%v\n",
		result2Str,
	)
	fmt.Printf(
		"Agent error 2:\n%v\n",
		err2Str,
	)
	fmt.Println("===================")
}
