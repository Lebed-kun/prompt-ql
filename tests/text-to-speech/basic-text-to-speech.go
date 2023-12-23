package texttospeechtest

import (
	"fmt"
	"path"

	interpretercore "gitlab.com/jbyte777/prompt-ql/v5/core"
	interpreter "gitlab.com/jbyte777/prompt-ql/v5/interpreter"
	testutils "gitlab.com/jbyte777/prompt-ql/v5/tests/utils"
)

// 28-09-2023: Works +++
func BasicTtsSyncQueryTest(
	openAiBaseUrl string,
	openAiKey string,
) {
	defaultGlobals := interpretercore.TGlobalVariablesTable{
		"save_to_file": testutils.SaveToFile,
		"file_path": path.Join("tests", "text-to-speech", "data", "speech_sync_query.mp3"),
	}
	interpreterInst := interpreter.New(
		interpreter.PromptQLOptions{
			OpenAiBaseUrl: openAiBaseUrl,
			OpenAiKey: openAiKey,
			DefaultExternalGlobals: defaultGlobals,
		},
	)

	result := interpreterInst.Instance.Execute(
		`
			=========^^^^^==========
			{~call fn=@save_to_file}
				{~get from=@file_path /}
				{~open_query_tts sync model="tts-1" voice="onyx"}
					The quick brown fox jumps over the lazy dog
				{/open_query_tts}
			{/call}
			=========+++++==========
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

// 28-09-2023: Works +++
func BasicTtsAsyncQueryTest(
	openAiBaseUrl string,
	openAiKey string,
) {
	defaultGlobals := interpretercore.TGlobalVariablesTable{
		"save_to_file": testutils.SaveToFile,
		"file_path": path.Join("tests", "text-to-speech", "data", "speech_async_query.mp3"),
	}
	interpreterInst := interpreter.New(
		interpreter.PromptQLOptions{
			OpenAiBaseUrl: openAiBaseUrl,
			OpenAiKey: openAiKey,
			DefaultExternalGlobals: defaultGlobals,
		},
	)

	result := interpreterInst.Instance.Execute(
		`
			=========^^^^^==========
			{~open_query_tts to="voice_res" model="tts-1" voice="onyx"}
				The quick brown fox jumps over the lazy dog
			{/open_query_tts}
			{~call fn=@save_to_file}
				{~get from=@file_path /}
				{~listen_query_tts from="voice_res" /}
			{/call}
			=========+++++==========
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

