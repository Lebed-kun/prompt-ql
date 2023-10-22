package interpretercore

import (
	stringsutils "gitlab.com/jbyte777/prompt-ql/v2/utils/strings"
)

type TFunctionInputChannel []interface{}

// [BEGIN] TODO: include this method to PromptQL v1.3.0 minor release
func (self TFunctionInputChannel) LatestCleanData() interface{} {
	if len(self) == 0 {
		return nil
	}
	
	for ptr := len(self) - 1; ptr >= 0; ptr -= 1 {
		data, isDataStr := self[ptr].(string)
		if !isDataStr && self[ptr] != nil {
			return self[ptr]
		}
		trimmedData := stringsutils.TrimWhitespace(data)
		if len(trimmedData) > 0 && trimmedData != " " {
			return trimmedData
		}
	}

	return nil
}
// [END] TODO: include this method to PromptQL v1.3.0 minor release
