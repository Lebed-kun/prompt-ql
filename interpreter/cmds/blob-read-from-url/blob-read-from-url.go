package blobreadfromurlcmd

import (
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
		_inputs interpreter.TFunctionInputChannelTable,
		_internalGlobals interpreter.TGlobalVariablesTable,
		_externalGlobals interpreter.TGlobalVariablesTable,
		execInfo interpreter.TExecutionInfo,
		_interpreter *interpreter.Interpreter,
	) interface{} {
		urlVar, err := getUrlVar(staticArgs, execInfo)
		if err != nil {
			return err
		}
	
		response, err := readFromUrl(urlVar, timeoutSec, execInfo)
		if err != nil {
			return err
		}

		return response
	}
}
