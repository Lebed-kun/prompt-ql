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

func New(
	openAiBaseUrl string,
	openAiKey string,
	openAiListenQueryTimeoutSec uint,
	customApisListenQueryTimeoutSec uint,
) *TPromptQL {
	apiInst := api.New(
		openAiBaseUrl,
		openAiKey,
		openAiListenQueryTimeoutSec,
	)
	customLLMApis := customapis.New(customApisListenQueryTimeoutSec)

	execFnTable := makeCmdTable(apiInst, customLLMApis)
	interpreterInst := interpreter.New(
		execFnTable,
		rootDataSwitch,
	)
	
	return &TPromptQL{
		Instance: interpreterInst,
		CustomApis: customLLMApis,
	}
}
