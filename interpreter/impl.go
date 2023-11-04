package interpreterimpl

import (
	interpreter "gitlab.com/jbyte777/prompt-ql/v3/core"
	api "gitlab.com/jbyte777/prompt-ql/v3/api"
	customapis "gitlab.com/jbyte777/prompt-ql/v3/custom-apis"
)

type TPromptQL struct {
	Instance *interpreter.Interpreter
	CustomApis *customapis.CustomModelsApis
}

type TPromptQLOptions struct {
	OpenAiBaseUrl string
	OpenAiKey string
	OpenAiListenQueryTimeoutSec uint
	CustomApisListenQueryTimeoutSec uint
	DefaultExternalGlobals interpreter.TGlobalVariablesTable
	DefaultExternalGlobalsMeta interpreter.TExternalGlobalsMetaTable
}

func New(options TPromptQLOptions) *TPromptQL {
	apiInst := api.New(
		options.OpenAiBaseUrl,
		options.OpenAiKey,
		options.OpenAiListenQueryTimeoutSec,
	)
	customModelsApis := customapis.New(options.CustomApisListenQueryTimeoutSec)

	execFnTable := makeCmdTable(apiInst, customModelsApis)
	interpreterInst := interpreter.New(
		execFnTable,
		rootDataSwitch,
		options.DefaultExternalGlobals,
		options.DefaultExternalGlobalsMeta,
	)
	
	return &TPromptQL{
		Instance: interpreterInst,
		CustomApis: customModelsApis,
	}
}
