package blobreadfromfilecmd

import (
	"fmt"
	"io/ioutil"
	"time"

	interpreter "gitlab.com/jbyte777/prompt-ql/v5/core"
)

func readFromFile(
	path string,
	timeoutSec uint,
	execInfo interpreter.TExecutionInfo,
) ([]byte, error) {
	resChan := make(chan []byte)
	errChan := make(chan error)

	go func() {
		file, err := ioutil.ReadFile(path)
		if err != nil {
			errChan <- fmt.Errorf(
				"!error (line=%v, char=%v): %v",
				execInfo.Line,
				execInfo.CharPos,
				err.Error(),
			)
		} else {
			resChan <- file
		}
	}()

	timeout := time.NewTimer(time.Second * time.Duration(timeoutSec))
	select {
	case res := <- resChan:
		return res, nil
	case err := <- errChan:
		return nil, err
	case <- timeout.C:
		return nil, fmt.Errorf(
			"!error (line=%v, char=%v): file read timeout",
			execInfo.Line,
			execInfo.CharPos,
		)
	}
}
