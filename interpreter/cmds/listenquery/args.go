package listenquerycmd

import (
	"errors"
	"fmt"

	interpreter "gitlab.com/jbyte777/prompt-ql/core"
)

func getFromVar(
	staticArgs interpreter.TFunctionArgumentsTable,
) (string, error) {
	var fromVar string
	rawFromVar, hasRawFromVar := staticArgs["from"]
	if !hasRawFromVar {
		return "", errors.New(
			"!error \"from\" parameter is required",
		)
	}
	var isFromVarStr bool
	fromVar, isFromVarStr = rawFromVar.(string)
	if !isFromVarStr {
		return "", fmt.Errorf(
			"!error \"from\" parameter is \"%v\" which is not string",
			rawFromVar,
		)
	}

	return fromVar, nil
}
