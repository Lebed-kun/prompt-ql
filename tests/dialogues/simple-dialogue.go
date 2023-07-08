package dialoguestests

import (
	"fmt"

	interpreter "gitlab.com/jbyte777/prompt-ql/interpreter"
)

func SimpleDialogTest(
	openAiBaseUrl string,
	openAiKey string,
) {
	interpreterInst := interpreter.New(
		openAiBaseUrl,
		openAiKey,
		40,
	)

	result := interpreterInst.Execute(
		`
			{~set to="reply1"}
				Hi, Bob! How are you?
			{/set}
			Alice: {~get from="reply1" /}
			{~open_query to="reply2" model="gpt-3.5-turbo-16k"}
				{~system}
					You are a good actor. Your name is Bob.
				{/system}
				{~get from="reply1" /}
			{/open_query}

			{~set to="reply2"}
				{~listen_query from="reply2" /}
			{/set}
			Bob: {~get from="reply2" /}
			{~open_query to="reply3" model="gpt-3.5-turbo-16k"}
				{~system}
					You are a good actor. Your name is Alice.
				{/system}
				{~assistant}
					{~get from="reply1" /}
				{/assistant}
				{~get from="reply2" /}
			{/open_query}

			{~set to="reply3"}
				{~listen_query from="reply3" /}
			{/set}
			Alice: {~get from="reply3" /}
			{~open_query to="reply4" model="gpt-3.5-turbo-16k"}
				{~system}
					You are a good actor. Your name is Bob.
				{/system}
				{~assistant}
					{~get from="reply2" /}
				{/assistant}
				{~get from="reply3" /}
			{/open_query}

			{~set to="reply4"}
				{~listen_query from="reply4" /}
			{/set}
			Bob: {~get from="reply4" /}
			{~open_query to="reply5" model="gpt-3.5-turbo-16k"}
				{~system}
					You are a good actor. Your name is Alice.
				{/system}
				{~assistant}
					{~get from="reply1" /}
				{/assistant}
				{~assistant}
					{~get from="reply3" /}
				{/assistant}
				{~get from="reply4" /}
			{/open_query}

			{~set to="reply5"}
				{~listen_query from="reply5" /}
			{/set}
			Alice: {~get from="reply5" /}
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
		"ChatGPT response:\n%v\n",
		resultStr,
	)
	fmt.Printf(
		"ChatGPT error:\n%v\n",
		errStr,
	)
	fmt.Println("===================")
}
