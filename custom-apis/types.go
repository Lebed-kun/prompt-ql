package customapis

import (
	interpreter "gitlab.com/jbyte777/prompt-ql/core"
)

type TDoQueryFuncTable map[string]TDoQueryFunc

type TDoQueryFunc func(
	model string,
	temperature float64,
	inputs interpreter.TFunctionInputChannelTable,
	execInfo interpreter.TExecutionInfo,
) ([]byte, error)

type TCustomQueryHandle struct {
	IsCustom bool
	ResultChan chan []byte
	ErrChan chan error
}
