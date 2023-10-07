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
		// [BEGIN] DONE: include this error text fix in patch v1.2.2
		return nil, fmt.Errorf(
			"!error (line=%v, char=%v): custom model named \"%v\" doesn't exist",
			execInfo.Line,
			execInfo.CharPos,
			model,
		)
		// [END] DONE: include this error text fix in patch v1.2.2
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

func (self *CustomLLMApis) GetAllModelsList() map[string]bool {
	res := make(map[string]bool, 0)

	for k := range self.llms {
		res[k] = true
	}

	return res
}
