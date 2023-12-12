package interpreterimpl

import (
	interpreter "gitlab.com/jbyte777/prompt-ql/v5/core"
	api "gitlab.com/jbyte777/prompt-ql/v5/api"
	customapis "gitlab.com/jbyte777/prompt-ql/v5/custom-apis"
	loggerapis "gitlab.com/jbyte777/prompt-ql/v5/logger-apis"
)

type PromptQL struct {
	Instance interpreter.InterpreterApi
	CustomApis customapis.CustomModelsMainApi
	LoggerApis loggerapis.LoggersMainApi
}

type PromptQLOptions struct {
	OpenAiBaseUrl string
	OpenAiKey string
	OpenAiListenQueryTimeoutSec uint
	CustomApisListenQueryTimeoutSec uint
	DefaultExternalGlobals interpreter.TGlobalVariablesTable
	DefaultExternalGlobalsMeta interpreter.TExternalGlobalsMetaTable
	PreinitializedInternalGlobals interpreter.TGlobalVariablesTable
	RestrictedCommands interpreter.TRestrictedCommands
	ReadFromFileTimeoutSec uint
	ReadFromUrlTimeoutSec uint
}

func New(options PromptQLOptions) *PromptQL {
	apiInst := api.New(
		options.OpenAiBaseUrl,
		options.OpenAiKey,
		options.OpenAiListenQueryTimeoutSec,
	)
	customModelsApis := customapis.New(options.CustomApisListenQueryTimeoutSec)
	loggerApis := loggerapis.New()

	execFnTable := makeCmdTable(
		apiInst,
		customModelsApis,
		loggerApis,
		options.ReadFromFileTimeoutSec,
		options.ReadFromUrlTimeoutSec,
	)
	
	interpreterOpts := interpreter.InterpreterConfig{
		ExecFnTable: execFnTable,
		DataSwitchFn: rootDataSwitch,
		DefaultExternalVars: options.DefaultExternalGlobals,
		DefaultExternalVarsMeta: options.DefaultExternalGlobalsMeta,
		RestrictedCommands: options.RestrictedCommands,
		CmdsMeta: cmdsMetaInfo,
		PreinitedInternalGlobals: options.PreinitializedInternalGlobals,
	}
	interpreterInst := interpreter.New(interpreterOpts)
	
	return &PromptQL{
		Instance: interpreterInst,
		CustomApis: customModelsApis,
		LoggerApis: loggerApis,
	}
}
