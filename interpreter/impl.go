package interpreterimpl

import (
	interpreter "gitlab.com/jbyte777/prompt-ql/core"
	api "gitlab.com/jbyte777/prompt-ql/api"
)

func New(
	openAiBaseUrl string,
	openAiKey string,
	listenQueryTimeoutSec uint,
) *interpreter.Interpreter {
	apiInst := api.New(
		openAiBaseUrl,
		openAiKey,
		listenQueryTimeoutSec,
	)

	execFnTable := makeCmdTable(apiInst)
	interpreterInst := interpreter.New(
		execFnTable,
		rootDataSwitch,
	)
	
	return interpreterInst
}
