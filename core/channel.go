package interpretercore

import (
	"fmt"
	"strings"

	stringsutils "gitlab.com/jbyte777/prompt-ql/v4/utils/strings"
)

type TFunctionInputChannel []interface{}

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

func (self TFunctionInputChannel) MergeIntoString() (string, error) {
	result := strings.Builder{}
	for _, arg := range self {
		argStr, isArgStr := arg.(string)

		if !isArgStr {
			return "", fmt.Errorf(
				"argument \"%v\" is not a string",
				arg,
			)
		}

		result.WriteString(argStr)
	}

	return stringsutils.TrimWhitespace(
		result.String(),
	), nil
}
