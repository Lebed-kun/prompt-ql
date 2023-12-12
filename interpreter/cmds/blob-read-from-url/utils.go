package blobreadfromurlcmd

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	interpreter "gitlab.com/jbyte777/prompt-ql/v5/core"
)

func readFromUrl(
	url string,
	timeoutSec uint,
	execInfo interpreter.TExecutionInfo,
) ([]byte, error) {
	resChan := make(chan []byte)
	errChan := make(chan error)

	go func() {
		client := http.Client{}
		request, err := http.NewRequest(
			"GET",
			url,
			nil,
		)
		if err != nil {
			errChan <- fmt.Errorf(
				"!error (line=%v, char=%v): %v",
				execInfo.Line,
				execInfo.CharPos,
				err.Error(),
			)
			return
		}

		response, err := client.Do(request)
		if err != nil {
			errChan <- fmt.Errorf(
				"!error (line=%v, char=%v): %v",
				execInfo.Line,
				execInfo.CharPos,
				err.Error(),
			)
			return
		}

		body, err := ioutil.ReadAll(response.Body)
		if err != nil {
			errChan <- fmt.Errorf(
				"!error (line=%v, char=%v): %v",
				execInfo.Line,
				execInfo.CharPos,
				err.Error(),
			)
			return
		}

		if response.StatusCode >= 400 {
			errChan <- fmt.Errorf(
				"!error (line=%v, char=%v): %v",
				execInfo.Line,
				execInfo.CharPos,
				string(body),
			)
			return
		}

		resChan <- body
	}()

	timeout := time.NewTimer(time.Second * time.Duration(timeoutSec))
	select {
	case res := <- resChan:
		return res, nil
	case err := <- errChan:
		return nil, err
	case <- timeout.C:
		return nil, fmt.Errorf(
			"!error (line=%v, char=%v): url read timeout",
			execInfo.Line,
			execInfo.CharPos,
		)
	}
}
