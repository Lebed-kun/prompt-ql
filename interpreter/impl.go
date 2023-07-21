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
	listenQueryTimeoutSec uint,
) *TPromptQL {
	apiInst := api.New(
		openAiBaseUrl,
		openAiKey,
		listenQueryTimeoutSec,
	)
	customLLMApis := customapis.New(listenQueryTimeoutSec)

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
