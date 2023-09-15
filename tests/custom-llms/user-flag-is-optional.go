package customllmstests

import (
	"fmt"
	interpreter "gitlab.com/jbyte777/prompt-ql/interpreter"
)

// Works ?
func UserFlagIsOptionalTest(
	pathToLlamaCommand string,
	pathToLlamaModel string,
) {
	interpreterInst := interpreter.New(
		"",
		"",
		0,
		400,
	)
	llamaDoQuery := makeLlamaDoQuery(pathToLlamaCommand, pathToLlamaModel)
	interpreterInst.CustomApis.RegisterLLMApi(
		"llama",
		llamaDoQuery,
	)

	result := interpreterInst.Instance.Execute(
		`
			{~open_query to="query1" model="llama"}
				{~system}
					You are a helpful assistant.
				{/system}
				I want a response to the following question:
				Tell me what is the Fourier Series
			{/open_query}
			{~listen_query from="query1" /}
		`,
		nil,
	)

	if result.Error != nil {
		panic(result.Error)
	}

	fmt.Println("===================")
	resultStr, _ := result.ResultDataStr()
	errStr, _ := result.ResultErrorStr()
	fmt.Printf(
		"Llama response:\n%v\n",
		resultStr,
	)
	fmt.Printf(
		"Llama error:\n%v\n",
		errStr,
	)
	fmt.Println("===================")
}
