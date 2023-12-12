package filetest

import (
	"fmt"

	interpretercore "gitlab.com/jbyte777/prompt-ql/v5/core"
	interpreter "gitlab.com/jbyte777/prompt-ql/v5/interpreter"
)

// 12-12-2023 Works +++
func BlobReadFromUrlTest() {
	defaultInternals := interpretercore.TGlobalVariablesTable{
		"_internal_fn": func(args []interface{}) interface{} {
			if len(args) < 1 {
				return fmt.Errorf("No byte slice is provided")
			}
			byteSlice, isByteSlice := args[0].([]byte)
			if !isByteSlice {
				return fmt.Errorf("No byte slice is provided")
			}

			return string(byteSlice)
		},
	}
	interpreterInst := interpreter.New(
		interpreter.PromptQLOptions{
			PreinitializedInternalGlobals: defaultInternals,
		},
	)

	result1 := interpreterInst.Instance.UnsafeExecute(
		`
			{~unsafe_preinit_vars /}
			{~set to="_url"}https://bvnf.space/{/set}
			{~call fn="_internal_fn"}
				{~blob_read_from_url url=$_url /}
			{/call}
		`,
	)

	if result1.Error != nil {
		panic(result1.Error)
	}

	fmt.Println("===================")
	result1Str, _ := result1.ResultDataStr()
	err1Str, _ := result1.ResultErrorStr()
	fmt.Printf(
		"Agent response 1:\n\"%v\"\n",
		result1Str,
	)
	fmt.Printf(
		"Agent error 1:\n\"%v\"\n",
		err1Str,
	)
	fmt.Println("===================")
}
