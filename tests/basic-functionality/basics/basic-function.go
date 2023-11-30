package basicfunctionalitytests

import (
	"fmt"

	interpretercore "gitlab.com/jbyte777/prompt-ql/v5/core"
	interpreter "gitlab.com/jbyte777/prompt-ql/v5/interpreter"
)

// Works ++++++
// 28-09-2023: Works on total regress +++
// 07-10-2023: Works +++
func BasicFunctionTest() {
	defaultGlobals := interpretercore.TGlobalVariablesTable{
		"postprocess": postProcessFunctionTest,
	}
	interpreterInst := interpreter.New(
		interpreter.TPromptQLOptions{
			DefaultExternalGlobals: defaultGlobals,
		},
	)

	result := interpreterInst.Instance.Execute(
		`
			{~set to="queryres"}
				1. Make a dish
				2. Eat it
				3. ...
			{/set}
			Raw result is:
			{~get from="queryres" /}
			==========++++++==========
			JSON result is:
			{~call fn=@postprocess }
				{~get from="queryres" /}
			{/call}
		`,
	)

	if result.Error != nil {
		panic(result.Error)
	}

	fmt.Println("===================")
	resultStr, _ := result.ResultDataStr()
	errStr, _ := result.ResultErrorStr()
	fmt.Printf(
		"ChatGPT response:\n%v\n",
		resultStr,
	)
	fmt.Printf(
		"ChatGPT error:\n%v\n",
		errStr,
	)
	fmt.Println("===================")
}
