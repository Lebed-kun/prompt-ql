package blobreadfromurlcmd

import (
	"fmt"
	interpreter "gitlab.com/jbyte777/prompt-ql/v5/core"
)

var defaultBlobReadUrlTimeoutSec uint = 20

func MakeBlobReadFromUrlCmd(
	blobReadUrlTimeoutSec uint,
) interpreter.TExecutedFunction {
	timeoutSec := blobReadUrlTimeoutSec
	if timeoutSec == 0 {
		timeoutSec = defaultBlobReadUrlTimeoutSec
	}

	return func(
		staticArgs interpreter.TFunctionArgumentsTable,
		inputs interpreter.TFunctionInputChannelTable,
		_internalGlobals interpreter.TGlobalVariablesTable,
		_externalGlobals interpreter.TGlobalVariablesTable,
		execInfo interpreter.TExecutionInfo,
		_interpreter *interpreter.Interpreter,
	) interface{} {
		urlVar, err := getUrlVar(staticArgs, execInfo)
		if err != nil {
			return err
		}
		methodVar, err := getMethodVar(staticArgs, execInfo)
		if err != nil {
			return err
		}
		
		var bodyStr string
		if methodVar != "GET" && methodVar != "OPTIONS" && inputs["data"] != nil {
			bodyStr, err = inputs["data"].MergeIntoString()
			if err != nil {
				return fmt.Errorf(
					"!error (line=%v, char=%v): %v",
					execInfo.Line,
					execInfo.CharPos,
					err.Error(),
				)
			}
		}
	
		response, err := readFromUrl(urlVar, methodVar, bodyStr, timeoutSec, execInfo)
		if err != nil {
			return err
		}

		return response
	}
}
