package loggerapis

import (
	interpreter "gitlab.com/jbyte777/prompt-ql/v5/core"
)

type TLoggerFuncTable map[string]TLoggerFunc

type TLoggerFunc func(
	execInfo interpreter.TExecutionInfo,
	inputs interpreter.TFunctionInputChannelTable,
	internalGlobals interpreter.TGlobalVariablesTable,
	externalGlobals interpreter.TGlobalVariablesTable,
) error
