package httputils

import (
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"time"
)

func DoHttpRequest(
	url string,
	method string,
	reqBody string,
	timeoutSec uint,
) ([]byte, error) {
	resChan := make(chan []byte)
	errChan := make(chan error)

	go func() {
		client := http.Client{}
		var request *http.Request
		var err error

		if len(reqBody) > 0 {
			bodyReqReader := strings.NewReader(reqBody)
			request, err = http.NewRequest(
				method,
				url,
				bodyReqReader,
			)
		} else {
			request, err = http.NewRequest(
				method,
				url,
				nil,
			)
		}
		
		if err != nil {
			errChan <- err
			return
		}

		response, err := client.Do(request)
		if err != nil {
			errChan <- err
			return
		}

		body, err := ioutil.ReadAll(response.Body)
		if err != nil {
			errChan <- err
			return
		}

		if response.StatusCode >= 400 {
			errChan <- fmt.Errorf("%v", string(body))
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
		return nil, errors.New("url read timeout")
	}
}
