package interpreterimpl

import (
	interpreter "gitlab.com/jbyte777/prompt-ql/core"
	api "gitlab.com/jbyte777/prompt-ql/api"
)

func New(
	openAiBaseUrl string,
	openAiKey string,
) *interpreter.Interpreter {
	apiInst := api.New(
		openAiBaseUrl,
		openAiKey,
	)

	execFnTable := makeCmdTable(apiInst)
	interpreterInst := interpreter.New(
		execFnTable,
		rootDataSwitch,
	)
	
	return interpreterInst
}
