package customapis

import (
	"errors"
	"fmt"
	"time"

	interpreter "gitlab.com/jbyte777/prompt-ql/v3/core"
	errorsutils "gitlab.com/jbyte777/prompt-ql/v3/utils/errors"
)

type CustomModelsApis struct {
	models                  TDoQueryFuncTable
	modelsMeta TCustomModelMetaTable
	listenQueryTimeoutSec uint
}

const defaultListenQueryTimeoutSec uint = 30

func New(listenQueryTimeoutSec uint) *CustomModelsApis {
	if listenQueryTimeoutSec == 0 {
		listenQueryTimeoutSec = defaultListenQueryTimeoutSec
	}

	return &CustomModelsApis{
		models:                  TDoQueryFuncTable{},
		modelsMeta: TCustomModelMetaTable{},
		listenQueryTimeoutSec: listenQueryTimeoutSec,
	}
}

func (self *CustomModelsApis) RegisterModelApi(
	name string,
	doQuery TDoQueryFunc,
	description string,
) {
	self.models[name] = doQuery
	
	if self.modelsMeta[name] == nil {
		self.modelsMeta[name] = &TCustomModelMetaInfo{
			Description: description,
		}
	}
	self.modelsMeta[name].Description = description
}

func (self *CustomModelsApis) UnregisterModelApi(name string) {
	delete(self.models, name)
	delete(self.modelsMeta, name)
}

func (self *CustomModelsApis) OpenQuery(
	model string,
	temperature float64,
	inputs interpreter.TFunctionInputChannelTable,
	execInfo interpreter.TExecutionInfo,
) (*TCustomQueryHandle, error) {
	doQuery, hasDoQuery := self.models[model]
	if !hasDoQuery {
		return nil, fmt.Errorf(
			"!error (line=%v, char=%v): custom model named \"%v\" doesn't exist",
			execInfo.Line,
			execInfo.CharPos,
			model,
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

func (self *CustomModelsApis) ListenQuery(
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
			"CustomModelsApis",
			"ListenQuery",
			errors.New("Timeout for listening query"),
		)
	}
}

func (self *CustomModelsApis) GetAllModelsList() map[string]string {
	res := make(map[string]string, 0)

	for k, v := range self.modelsMeta {
		res[k] = v.Description
	}

	return res
}
