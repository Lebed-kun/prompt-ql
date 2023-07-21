package customapis

import (
	"errors"
	"fmt"
	"time"

	errorsutils "gitlab.com/jbyte777/prompt-ql/utils/errors"
	interpreter "gitlab.com/jbyte777/prompt-ql/core"
)

type CustomLLMApis struct {
	llms TDoQueryFuncTable
	listenQueryTimeoutSec uint
}

const defaultListenQueryTimeoutSec uint = 30

func New(listenQueryTimeoutSec uint) *CustomLLMApis {
	if listenQueryTimeoutSec == 0 {
		listenQueryTimeoutSec = defaultListenQueryTimeoutSec
	}
	
	return &CustomLLMApis{
		llms: TDoQueryFuncTable{},
		listenQueryTimeoutSec: listenQueryTimeoutSec,
	}
}

func (self *CustomLLMApis) RegisterLLMApi(
	name string,
	doQuery TDoQueryFunc,
) {
	self.llms[name] = doQuery
}

func (self *CustomLLMApis) OpenQuery(
	model string,
	temperature float64,
	inputs interpreter.TFunctionInputChannelTable,
	execInfo interpreter.TExecutionInfo,
) *TCustomQueryHandle {
	resChan := make(chan []byte)
	errChan := make(chan error)

	go func() {
		doQuery, hasDoQuery := self.llms[model]
		if !hasDoQuery {
			errChan <- fmt.Errorf(
				"!error (line=%v, char=%v): prompts are empty",
				execInfo.Line,
				execInfo.CharPos,
			)
			return
		}

		res, err := doQuery(
			model,
			temperature,
			inputs,
			execInfo,
		)

		if err != nil {
			errChan <- err
		} else {
			resChan <- res
		}
	}()

	return &TCustomQueryHandle{
		IsCustom: true,
		ResultChan: resChan,
		ErrChan: errChan,
	}
}

func (self *CustomLLMApis) ListenQuery(
	queryHandle *TCustomQueryHandle,
) ([]byte, error) {
	timer := time.NewTimer(
		time.Second * time.Duration(self.listenQueryTimeoutSec),
	)

	select {
	case res := <-queryHandle.ResultChan:
		return res, nil
	case err := <-queryHandle.ErrChan:
		return nil, err
	case <-timer.C:
		return nil, errorsutils.LogError(
			"CustomLLMApis",
			"ListenQuery",
			errors.New("Timeout for listening query"),
		)
	}
}
