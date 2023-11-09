package customapis

import (
	interpreter "gitlab.com/jbyte777/prompt-ql/v4/core"
)

type TDoQueryFuncTable map[string]TDoQueryFunc

type TDoQueryFunc func(
	model string,
	temperature float64,
	inputs interpreter.TFunctionInputChannelTable,
	execInfo interpreter.TExecutionInfo,
) (string, error)

type TCustomModelMetaTable map[string]*TCustomModelMetaInfo

type TCustomModelMetaInfo struct {
	Description string
}

type TCustomQueryHandle struct {
	IsCustom bool
	ResultChan chan string
	ErrChan chan error
}
