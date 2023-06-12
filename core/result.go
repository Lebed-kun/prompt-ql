package interpretercore

import (
	"strings"
	stringsutils "gitlab.com/jbyte777/prompt-ql/utils/strings"
)

type TInterpreterResult struct {
	Result TFunctionInputChannelTable
	Error    error
	Finished bool
}

func (self *TInterpreterResult) stringifyResultChan(chanName string) (string, bool) {
	dataChan, hasDataChan := self.Result[chanName]
	if !hasDataChan {
		return "", false
	}
	
	for _, arg := range dataChan {
		_, isArgStr := arg.(string)
		if !isArgStr {
			return "", false
		}
	}

	result := strings.Builder{}
	for _, arg := range dataChan {
		argStr := arg.(string)
		result.WriteString(argStr)
	}

	return stringsutils.TrimWhitespace(
		result.String(),
	), true
}

func (self *TInterpreterResult) ResultDataStr() (string, bool) {
	return self.stringifyResultChan("data")
}

func (self *TInterpreterResult) ResultErrorStr() (string, bool) {
	return self.stringifyResultChan("error")
}
