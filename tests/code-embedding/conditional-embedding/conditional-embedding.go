package conditionalembeddingtest

import (
	"fmt"

	interpretercore "gitlab.com/jbyte777/prompt-ql/v5/core"
	interpreter "gitlab.com/jbyte777/prompt-ql/v5/interpreter"
)

// 04-11-2023: Works +++
func ConditionalEmbeddingTest() {
	defaultGlobals := interpretercore.TGlobalVariablesTable{
		"myFunc": func(args []interface{}) interface{} {
			argStr, _ := args[0].(string)
			return argStr
		},
		"myCond": func(args []interface{}) bool {
			argStr1, _ := args[0].(string)
			argStr2, _ := args[1].(string)
			return len(argStr1) > len(argStr2)
		},
	}
	interpreterInst := interpreter.New(
		interpreter.TPromptQLOptions{
			DefaultExternalGlobals: defaultGlobals,
		},
	)

	result := interpreterInst.Instance.Execute(
		`
			{~call fn=@myFunc}{~data}The 1st line of PromptQL will be executed{/data}{/call}
			{~embed_if cond=@myCond}
				{~data}Hello, world!{/data}
				{~data}kek{/data}
				<%
				  {~call fn=@myFunc}{~data}The 2nd line of PromptQL will be embedded{/data}{/call}
				%>
				<%
				  {~call fn=@myFunc}{~data}The 3rd line of PromptQL won't be embedded{/data}{/call}
				%>
			{/embed_if}
			{~call fn=@myFunc}{~data}The 4th line of PromptQL will be executed{/data}{/call}
			{~embed_if cond=@myCond}
				{~data}sos{/data}
				{~data}Hello{/data}
				<%
				  {~call fn=@myFunc}{~data}The 5th line of PromptQL won't be embedded{/data}{/call}
				%>
				<%
				  {~call fn=@myFunc}{~data}The 6th line of PromptQL will be embedded{/data}{/call}
				%>
			{/embed_if}
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
