package interpreterimpl

import (
	interpreter "gitlab.com/jbyte777/prompt-ql/v4/core"
	api "gitlab.com/jbyte777/prompt-ql/v4/api"
	customapis "gitlab.com/jbyte777/prompt-ql/v4/custom-apis"
)

type PromptQL struct {
	Instance interpreter.InterpreterApi
	CustomApis customapis.CustomModelsMainApi
}

type TPromptQLOptions struct {
	OpenAiBaseUrl string
	OpenAiKey string
	OpenAiListenQueryTimeoutSec uint
	CustomApisListenQueryTimeoutSec uint
	DefaultExternalGlobals interpreter.TGlobalVariablesTable
	DefaultExternalGlobalsMeta interpreter.TExternalGlobalsMetaTable
	RestrictedCommands interpreter.TRestrictedCommands
}

func New(options TPromptQLOptions) *PromptQL {
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
		options.RestrictedCommands,
	)
	
	return &PromptQL{
		Instance: interpreterInst,
		CustomApis: customModelsApis,
	}
}
