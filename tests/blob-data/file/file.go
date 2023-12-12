package filetest

import (
	"fmt"
	"path"

	interpretercore "gitlab.com/jbyte777/prompt-ql/v5/core"
	interpreter "gitlab.com/jbyte777/prompt-ql/v5/interpreter"
)

// 12-12-2023 Works +++
func BlobReadFromFileTest() {
	defaultInternals := interpretercore.TGlobalVariablesTable{
		"_internal_fn": func(args []interface{}) interface{} {
			if len(args) < 1 {
				return fmt.Errorf("No byte slice is provided")
			}
			byteSlice, isByteSlice := args[0].([]byte)
			if !isByteSlice {
				return fmt.Errorf("Argument is not a byte slice")
			}

			return string(byteSlice)
		},
		"_file_path": path.Join(
			"tests", "blob-data", "file", "test-file.txt",
		),
	}
	interpreterInst := interpreter.New(
		interpreter.PromptQLOptions{
			PreinitializedInternalGlobals: defaultInternals,
		},
	)

	result1 := interpreterInst.Instance.UnsafeExecute(
		`
			{~unsafe_preinit_vars /}
			{~call fn="_internal_fn"}
				{~blob_from_file path=$_file_path /}
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
		"Agent response 1:\n%v\n",
		result1Str,
	)
	fmt.Printf(
		"Agent error 1:\n%v\n",
		err1Str,
	)
	fmt.Println("===================")
}
