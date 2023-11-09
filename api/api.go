package api

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	errorsutils "gitlab.com/jbyte777/prompt-ql/v4/utils/errors"
)

type GptApi struct {
	openAiBaseUrl         string
	openAiKey             string
	listenQueryTimeoutSec uint
}

const defaultListenQueryTimeoutSec uint = 30

func New(
	openAiBaseUrl string,
	openAiKey string,
	listenQueryTimeoutSec uint,
) *GptApi {
	if listenQueryTimeoutSec == 0 {
		listenQueryTimeoutSec = defaultListenQueryTimeoutSec
	}

	return &GptApi{
		openAiBaseUrl:         openAiBaseUrl,
		openAiKey:             fmt.Sprintf("Bearer %v", openAiKey),
		listenQueryTimeoutSec: listenQueryTimeoutSec,
	}
}

func (self *GptApi) doQuery(
	model string,
	temperature float64,
	prompts []TMessage,
) (*TGptApiResponse, error) {
	client := &http.Client{}

	requestBody := TGptApiRequest{
		Model:       model,
		Messages:    prompts,
		Temperature: temperature,
		N:           1,
	}
	requestBodyBytes, _ := json.Marshal(requestBody)
	requestBodyBuff := bytes.NewBuffer(requestBodyBytes)

	reqUrl := fmt.Sprintf(
		"%v/v1/chat/completions",
		self.openAiBaseUrl,
	)
	request, _ := http.NewRequest(
		"POST",
		reqUrl,
		requestBodyBuff,
	)
	request.Header.Add("Authorization", self.openAiKey)
	request.Header.Add("Content-Type", "application/json")

	response, err := client.Do(request)
	if err != nil {
		return nil, errorsutils.LogError(
			"GptApi",
			"doQuery",
			err,
		)
	}
	if response.StatusCode != 200 {
		resBody, _ := ioutil.ReadAll(response.Body)

		return nil, errorsutils.LogError(
			"GptApi",
			"doQuery",
			errors.New(string(resBody)),
		)
	}

	rawResBody, _ := ioutil.ReadAll(response.Body)
	var resBody TGptApiResponse
	err = json.Unmarshal(rawResBody, &resBody)
	if err != nil {
		return nil, errorsutils.LogError(
			"GptApi",
			"doQuery",
			err,
		)
	}

	return &resBody, nil
}

func (self *GptApi) OpenQuery(
	model string,
	temperature float64,
	prompts []TMessage,
) (*TQueryHandle, error) {
	_, isModelSupported := supportedOpenAiModels[model]
	if !isModelSupported {
		return nil, fmt.Errorf(
			"model \"%v\" is not supported by OpenAI",
			model,
		)
	}

	resChan := make(chan *TGptApiResponse)
	errChan := make(chan error)

	go func() {
		res, err := self.doQuery(
			model,
			temperature,
			prompts,
		)

		if err != nil {
			errChan <- err
		} else {
			resChan <- res
		}
	}()

	return &TQueryHandle{
		ResultChan: resChan,
		ErrChan:    errChan,
	}, nil
}

func (self *GptApi) ListenQuery(
	queryHandle *TQueryHandle,
) (*TGptApiResponse, error) {
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
			"GptApi",
			"ListenQuery",
			errors.New("Timeout for listening query"),
		)
	}
}

func (self *GptApi) IsModelSupported(model string) bool {
	_, isModelSupported := supportedOpenAiModels[model]
	return isModelSupported
}

func (self *GptApi) GetAllModelsList() map[string]string {
	res := make(map[string]string, 0)

	for k, v := range supportedOpenAiModels {
		res[k] = v
	}

	return res
}
