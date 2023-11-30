package interpreterimpl

import (
	interpreter "gitlab.com/jbyte777/prompt-ql/v5/core"
	api "gitlab.com/jbyte777/prompt-ql/v5/api"
	customapis "gitlab.com/jbyte777/prompt-ql/v5/custom-apis"
	loggerapis "gitlab.com/jbyte777/prompt-ql/v5/logger-apis"
)

type TPromptQL struct {
	Instance *interpreter.Interpreter
	CustomApis *customapis.CustomModelsApis
	LoggerApis *loggerapis.LoggerApis
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

func New(options TPromptQLOptions) *TPromptQL {
	apiInst := api.New(
		options.OpenAiBaseUrl,
		options.OpenAiKey,
		options.OpenAiListenQueryTimeoutSec,
	)
	customModelsApis := customapis.New(options.CustomApisListenQueryTimeoutSec)
	loggerApis := loggerapis.New()

	execFnTable := makeCmdTable(apiInst, customModelsApis, loggerApis)
	
	interpreterInst := interpreter.New(
		execFnTable,
		rootDataSwitch,
		options.DefaultExternalGlobals,
		options.DefaultExternalGlobalsMeta,
		options.RestrictedCommands,
		cmdsMetaInfo,
	)
	
	return &TPromptQL{
		Instance: interpreterInst,
		CustomApis: customModelsApis,
		LoggerApis: loggerApis,
	}
}
