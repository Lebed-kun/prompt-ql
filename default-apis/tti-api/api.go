package ttiapi

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	errorsutils "gitlab.com/jbyte777/prompt-ql/v5/utils/errors"
)

type TtiApi struct {
	openAiBaseUrl         string
	openAiKey             string
	listenQueryTimeoutSec uint
}

const defaultListenQueryTimeoutSec uint = 50

func New(
	openAiBaseUrl string,
	openAiKey string,
	listenQueryTimeoutSec uint,
) *TtiApi {
	if len(openAiBaseUrl) == 0 {
		openAiBaseUrl = "https://api.openai.com"
	}

	if listenQueryTimeoutSec == 0 {
		listenQueryTimeoutSec = defaultListenQueryTimeoutSec
	}

	return &TtiApi{
		openAiBaseUrl:         openAiBaseUrl,
		openAiKey:             fmt.Sprintf("Bearer %v", openAiKey),
		listenQueryTimeoutSec: listenQueryTimeoutSec,
	}
}

func (self *TtiApi) doQuery(
	model string,
	prompt string,
	width uint,
	height uint,
	responseFormat string,
) (*TTtiApiResponse, error) {
	client := &http.Client{}

	requestBody := TTtiApiRequest{
		Model:          model,
		Prompt:         prompt,
		N:              1,
		Size:           fmt.Sprintf("%vx%v", width, height),
		ResponseFormat: responseFormat,
	}
	requestBodyBytes, _ := json.Marshal(requestBody)
	requestBodyBuff := bytes.NewBuffer(requestBodyBytes)

	reqUrl := fmt.Sprintf(
		"%v/v1/audio/speech",
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
			"TtiApi",
			"doQuery",
			err,
		)
	}
	if response.StatusCode != 200 {
		resBody, _ := ioutil.ReadAll(response.Body)

		return nil, errorsutils.LogError(
			"TtiApi",
			"doQuery",
			errors.New(string(resBody)),
		)
	}

	rawResBody, _ := ioutil.ReadAll(response.Body)
	var resBody TTtiApiResponse
	err = json.Unmarshal(rawResBody, &resBody)
	if err != nil {
		return nil, errorsutils.LogError(
			"TtiApi",
			"doQuery",
			err,
		)
	}

	return &resBody, nil
}

func (self *TtiApi) OpenQuery(
	model string,
	prompt string,
	width uint,
	height uint,
	responseFormat string,
) (*TQueryHandle, error) {
	_, isModelSupported := supportedOpenAiModels[model]
	if !isModelSupported {
		return nil, fmt.Errorf(
			"model \"%v\" is not supported by OpenAI",
			model,
		)
	}

	resChan := make(chan *TTtiApiResponse)
	errChan := make(chan error)

	go func() {
		res, err := self.doQuery(
			model,
			prompt,
			width,
			height,
			responseFormat,
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

func (self *TtiApi) ListenQuery(
	queryHandle *TQueryHandle,
) (*TTtiApiResponse, error) {
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
			"TtiApi",
			"ListenQuery",
			errors.New("Timeout for listening query"),
		)
	}
}

func (self *TtiApi) IsModelSupported(model string) bool {
	_, isModelSupported := supportedOpenAiModels[model]
	return isModelSupported
}

func (self *TtiApi) GetAllModelsList() map[string]string {
	res := make(map[string]string, 0)

	for k, v := range supportedOpenAiModels {
		res[k] = v
	}

	return res
}
