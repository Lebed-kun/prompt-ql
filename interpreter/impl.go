package interpreterimpl

import (
	interpreter "gitlab.com/jbyte777/prompt-ql/core"
	api "gitlab.com/jbyte777/prompt-ql/api"
)

func NewInterpreter(
	openAiBaseUrl string,
	openAiKey string,
) *interpreter.PromptQLInterpreter {
	apiInst := api.New(
		openAiBaseUrl,
		openAiKey,
	)

	execFnTable := makeCmdTable(apiInst)
	interpreterInst := interpreter.New(
		execFnTable,
		// TODO: add switch fn according to specsc =^_^=
		nil,
	)
	
	return interpreterInst
}
