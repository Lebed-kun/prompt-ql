package blobreadfromfilecmd

import (
	interpreter "gitlab.com/jbyte777/prompt-ql/v5/core"
)

var defaultBlobReadFileTimeoutSec uint = 3

func MakeBlobReadFromFileCmd(
	blobReadFileTimeoutSec uint,
) interpreter.TExecutedFunction {
	timeoutSec := blobReadFileTimeoutSec
	if timeoutSec == 0 {
		timeoutSec = defaultBlobReadFileTimeoutSec
	}

	return func(
		staticArgs interpreter.TFunctionArgumentsTable,
		_inputs interpreter.TFunctionInputChannelTable,
		_internalGlobals interpreter.TGlobalVariablesTable,
		_externalGlobals interpreter.TGlobalVariablesTable,
		execInfo interpreter.TExecutionInfo,
		_interpreter *interpreter.Interpreter,
	) interface{} {
		pathVar, err := getPathVar(staticArgs, execInfo)
		if err != nil {
			return err
		}
	
		file, err := readFromFile(pathVar, timeoutSec, execInfo)
		if err != nil {
			return err
		}

		return file
	}
}
