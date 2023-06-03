package api

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	errorsutils "gitlab.com/jbyte777/prompt-ql/utils/errors"
)

type GptApi struct {
	openAiBaseUrl string
	openAiKey     string
}

const ListenQueryTimeoutSec int = 15

func New(
	openAiBaseUrl string,
	openAiKey string,
) *GptApi {
	return &GptApi{
		openAiBaseUrl: openAiBaseUrl,
		openAiKey:     fmt.Sprintf("Bearer %v", openAiKey),
	}
}

func (self *GptApi) doQuery(
	model string,
	temperature float64,
	n int,
	prompts []TMessage,
) (*TGptApiResponse, error) {
	client := &http.Client{}

	requestBody := TGptApiRequest{
		Model:       model,
		Messages:    prompts,
		Temperature: temperature,
		N:           n,
	}
	requestBodyStr, _ := json.Marshal(requestBody)
	requestBodyBytes := bytes.NewBuffer(requestBodyStr)

	reqUrl := fmt.Sprintf(
		"%v/v1/chat/completions",
		self.openAiBaseUrl,
	)
	request, _ := http.NewRequest(
		"POST",
		reqUrl,
		requestBodyBytes,
	)
	request.Header.Add("Authorization", self.openAiKey)

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
	n int,
	prompts []TMessage,
) (chan *TGptApiResponse, chan error) {
	resChan := make(chan *TGptApiResponse)
	errChan := make(chan error)

	go func() {
		res, err := self.doQuery(
			model,
			temperature,
			n,
			prompts,
		)

		if err != nil {
			errChan <- err
		} else {
			resChan <- res
		}
	}()

	return resChan, errChan
}

func (self *GptApi) ListenQuery(
	resChan chan *TGptApiResponse,
	errChan chan error,
) (*TGptApiResponse, error) {
	timer := time.NewTimer(
		time.Second * time.Duration(ListenQueryTimeoutSec),
	)

	select {
	case res := <-resChan:
		return res, nil
	case err := <-errChan:
		return nil, err
	case <-timer.C:
		return nil, errorsutils.LogError(
			"GptApi",
			"ListenQuery",
			errors.New("Timeout for listening query"),
		)
	}
}
