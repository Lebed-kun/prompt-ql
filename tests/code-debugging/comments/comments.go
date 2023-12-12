package commentstests

import (
	"fmt"

	interpretercore "gitlab.com/jbyte777/prompt-ql/v5/core"
	interpreter "gitlab.com/jbyte777/prompt-ql/v5/interpreter"
)

// 30-11-2023: Works +++
func BasicCommentTest() {
	defaultGlobals := interpretercore.TGlobalVariablesTable{
		"myFunc": func(args []interface{}) interface{} {
			argStr, _ := args[0].(string)
			return argStr
		},
	}
	interpreterInst := interpreter.New(
		interpreter.PromptQLOptions{
			DefaultExternalGlobals: defaultGlobals,
		},
	)

	result := interpreterInst.Instance.Execute(
		`
			{~call fn=@myFunc}{~data}The first line of PromptQL will be executed{/data}{/call}
			<%
				{~call fn=@myFunc}{~data}The second line of PromptQL will be just printed to result{/data}{/call}
			%>
			{~call fn=@myFunc}{~data}The third line of PromptQL will be executed{/data}{/call}
			<~
				{~call fn=@myFunc}{~data}The fourth line of PromptQL will be just ignored{/data}{/call}
			~>
			{~call fn=@myFunc}{~data}The fifth line of PromptQL will be executed{/data}{/call}
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
