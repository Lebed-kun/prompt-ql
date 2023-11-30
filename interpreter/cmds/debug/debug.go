package debugcmd

import (
	"fmt"
	interpreter "gitlab.com/jbyte777/prompt-ql/v5/core"
	loggerapismod "gitlab.com/jbyte777/prompt-ql/v5/logger-apis"
)

func MakeDebugCmd(
	loggerApis *loggerapismod.LoggerApis,
) interpreter.TExecutedFunction {
	return func(
		staticArgs interpreter.TFunctionArgumentsTable,
		inputs interpreter.TFunctionInputChannelTable,
		internalGlobals interpreter.TGlobalVariablesTable,
		externalGlobals interpreter.TGlobalVariablesTable,
		execInfo interpreter.TExecutionInfo,
		_interpreter *interpreter.Interpreter,
	) interface{} {
		loggerVar, err := getLoggerVar(staticArgs, execInfo)
		if err != nil {
			return err
		}
		
		err = loggerApis.Log(
			loggerVar,
			execInfo,
			inputs,
			internalGlobals,
			externalGlobals,
		)
		if err != nil {
			fmt.Printf(
				"ERROR LOGGING (line=%v, char=%v): %v\n",
				execInfo.Line,
				execInfo.CharPos,
				err.Error(),
			)
			loggerApis.Log(
				"default",
				execInfo,
				inputs,
				internalGlobals,
				externalGlobals,
			)
			return nil
		}

		return nil
	}
}
