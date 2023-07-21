package customapis

import (
	"errors"
	"fmt"
	"time"

	interpreter "gitlab.com/jbyte777/prompt-ql/core"
	errorsutils "gitlab.com/jbyte777/prompt-ql/utils/errors"
)

type CustomLLMApis struct {
	llms                  TDoQueryFuncTable
	listenQueryTimeoutSec uint
}

const defaultListenQueryTimeoutSec uint = 30

func New(listenQueryTimeoutSec uint) *CustomLLMApis {
	if listenQueryTimeoutSec == 0 {
		listenQueryTimeoutSec = defaultListenQueryTimeoutSec
	}

	return &CustomLLMApis{
		llms:                  TDoQueryFuncTable{},
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
) (*TCustomQueryHandle, error) {
	doQuery, hasDoQuery := self.llms[model]
	if !hasDoQuery {
		return nil, fmt.Errorf(
			"!error (line=%v, char=%v): prompts are empty",
			execInfo.Line,
			execInfo.CharPos,
		)
	}

	resChan := make(chan string)
	errChan := make(chan error)

	go func() {
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
		IsCustom:   true,
		ResultChan: resChan,
		ErrChan:    errChan,
	}, nil
}

func (self *CustomLLMApis) ListenQuery(
	queryHandle *TCustomQueryHandle,
) (string, error) {
	timer := time.NewTimer(
		time.Second * time.Duration(self.listenQueryTimeoutSec),
	)

	select {
	case res := <-queryHandle.ResultChan:
		return res, nil
	case err := <-queryHandle.ErrChan:
		return "", err
	case <-timer.C:
		return "", errorsutils.LogError(
			"CustomLLMApis",
			"ListenQuery",
			errors.New("Timeout for listening query"),
		)
	}
}
