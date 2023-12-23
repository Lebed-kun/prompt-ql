package blobreadfromurlcmd

import (
	"fmt"
	httputils "gitlab.com/jbyte777/prompt-ql/v5/utils/http"
	interpreter "gitlab.com/jbyte777/prompt-ql/v5/core"
)

func readFromUrl(
	url string,
	method string,
	bodyStr string,
	timeoutSec uint,
	execInfo interpreter.TExecutionInfo,
) ([]byte, error) {
	res, err := httputils.DoHttpRequest(
		url,
		method,
		bodyStr,
		timeoutSec,
	)
	if err != nil {
		return nil, fmt.Errorf(
			"!error (line=%v, char=%v): %v",
			execInfo.Line,
			execInfo.CharPos,
			err.Error(),
		)
	}

	return res, nil
}
