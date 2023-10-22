package interpretercore

import (
	"strings"
	stringsutils "gitlab.com/jbyte777/prompt-ql/v2/utils/strings"
	promptmsgutils "gitlab.com/jbyte777/prompt-ql/v2/utils/promptmsg"
)

type TInterpreterResult struct {
	Result TFunctionInputChannelTable
	Error    error
	Complete bool
}

func (self *TInterpreterResult) ResultDataStr() (string, bool) {
	dataChan, hasDataChan := self.Result["data"]
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

func (self *TInterpreterResult) ResultLatestData(chanName string) interface{} {
	dataChan, hasDataChan := self.Result[chanName]
	if !hasDataChan {
		return nil
	}
	if len(dataChan) == 0 {
		return nil
	}
	
	for ptr := len(dataChan) - 1; ptr >= 0; ptr -= 1 {
		data, isDataStr := dataChan[ptr].(string)
		if !isDataStr && dataChan[ptr] != nil {
			return dataChan[ptr]
		}
		trimmedData := stringsutils.TrimWhitespace(data)
		if len(trimmedData) > 0 && trimmedData != " " {
			return trimmedData
		}
	}

	return nil
}

func (self *TInterpreterResult) ResultErrorStr() (string, bool) {
	errChan, hasErrChan := self.Result["error"]
	if !hasErrChan {
		return "", false
	}

	for _, arg := range errChan {
		_, isArgStr := arg.(string)
		_, isArgErr := arg.(error)

		if !isArgStr && !isArgErr {
			return "", false
		}
	}

	result := strings.Builder{}
	for _, arg := range errChan {
		result.WriteString("ERROR: ")

		argStr, isArgStr := arg.(string)
		if isArgStr {
			result.WriteString(
				promptmsgutils.ReplacePromptMsgPrefix(
					argStr,
					"",
				),
			)
		} else {
			argErr := arg.(error)
			result.WriteString(
				promptmsgutils.ReplacePromptMsgPrefix(
					argErr.Error(),
					"",
				),
			)
		}

		result.WriteString(";\n")
	}

	return stringsutils.TrimWhitespace(
		result.String(),
	), true
}
