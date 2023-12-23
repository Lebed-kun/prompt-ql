package texttoimagetest

import (
	"fmt"
	"path"

	interpretercore "gitlab.com/jbyte777/prompt-ql/v5/core"
	interpreter "gitlab.com/jbyte777/prompt-ql/v5/interpreter"
	testutils "gitlab.com/jbyte777/prompt-ql/v5/tests/utils"
)

// 23-12-2023: Works ???
func BasicTtiSyncQueryTest(
	openAiBaseUrl string,
	openAiKey string,
) {
	defaultGlobals := interpretercore.TGlobalVariablesTable{
		"save_to_file": testutils.SaveToFile,
		"file_path": path.Join("tests", "text-to-image", "data", "image_sync_query.webp"),
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
				{~open_query_tti sync model="dall-e-2" width="720" height="720" response_format="url"}
					The quick brown fox jumps over the lazy dog
				{/open_query_tti}
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

// 23-12-2023: Works ???
func BasicTtsAsyncQueryTest(
	openAiBaseUrl string,
	openAiKey string,
) {
	defaultGlobals := interpretercore.TGlobalVariablesTable{
		"save_to_file": testutils.SaveToFile,
		"file_path": path.Join("tests", "text-to-speech", "data", "image_async_query.webp"),
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
			{~open_query_tti to="img_res" model="dall-e-2" width="720" height="720" response_format="url"}
					The quick brown fox jumps over the lazy dog
				{/open_query_tti}
			{~call fn=@save_to_file}
				{~get from=@file_path /}
				{~listen_query_tti from="img_res" /}
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

// 23-12-2023: Works ???
func BasicTtiSyncBase64QueryTest(
	openAiBaseUrl string,
	openAiKey string,
) {
	defaultGlobals := interpretercore.TGlobalVariablesTable{
		"save_to_file": testutils.SaveToFile,
		"file_path": path.Join("tests", "text-to-image", "data", "image_sync_query_b64.webp"),
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
				{~open_query_tti sync model="dall-e-2" width="720" height="720" response_format="b64_json"}
					The quick brown fox jumps over the lazy dog
				{/open_query_tti}
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

