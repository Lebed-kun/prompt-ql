package debugtest

import (
	"fmt"

	interpretercore "gitlab.com/jbyte777/prompt-ql/v5/core"
	interpreter "gitlab.com/jbyte777/prompt-ql/v5/interpreter"
)

// 30-11-2023: Works +++
func BasicDebugTest() {
	defaultGlobals := interpretercore.TGlobalVariablesTable{
		"myFunc": func(args []interface{}) interface{} {
			argStr, _ := args[0].(string)
			return argStr
		},
		"myVar": "Just a global var",
	}
	interpreterInst := interpreter.New(
		interpreter.PromptQLOptions{
			DefaultExternalGlobals: defaultGlobals,
		},
	)

	result := interpreterInst.Instance.Execute(
		`
			{~debug}
				{~error}Some error line{/error}
				{~set to="myVar"}MyVar value{/set}
				{~set to="myVar2"}MyVar2 value{/set}
				{~user}1st user line{/user}
				{~system}1st system line{/system}
				{~user}2nd user line{/user}
				{~assistant}1st assistant line{/assistant}
				{~system}2nd system line{/system}
				{~assistant}2nd assistant line{/assistant}
				{~data}1st data line{/data}
				{~get from="myVar" /}
				{~get from="myVar2" /}
			{/debug}
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
