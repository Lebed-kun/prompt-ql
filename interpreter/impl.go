package interpreterimpl

import (
	interpreter "gitlab.com/jbyte777/prompt-ql/core"
	api "gitlab.com/jbyte777/prompt-ql/api"
	customapis "gitlab.com/jbyte777/prompt-ql/custom-apis"
)

type TPromptQL struct {
	Instance *interpreter.Interpreter
	CustomApis *customapis.CustomLLMApis
}

type TPromptQLOptions struct {
	OpenAiBaseUrl string
	OpenAiKey string
	OpenAiListenQueryTimeoutSec uint
	CustomApisListenQueryTimeoutSec uint
	DefaultExternalGlobals interpreter.TGlobalVariablesTable
}

func New(options TPromptQLOptions) *TPromptQL {
	apiInst := api.New(
		options.OpenAiBaseUrl,
		options.OpenAiKey,
		options.OpenAiListenQueryTimeoutSec,
	)
	customLLMApis := customapis.New(options.CustomApisListenQueryTimeoutSec)

	execFnTable := makeCmdTable(apiInst, customLLMApis)
	interpreterInst := interpreter.New(
		execFnTable,
		rootDataSwitch,
		options.DefaultExternalGlobals,
	)
	
	return &TPromptQL{
		Instance: interpreterInst,
		CustomApis: customLLMApis,
	}
}
