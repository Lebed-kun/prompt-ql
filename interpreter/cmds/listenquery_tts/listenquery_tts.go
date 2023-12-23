package listenqueryttscmd

import (
	"fmt"
	ttsapi "gitlab.com/jbyte777/prompt-ql/v5/default-apis/tts-api"
	interpreter "gitlab.com/jbyte777/prompt-ql/v5/core"
)

func MakeListenQueryTtsCmd(
	ttsApi *ttsapi.TtsApi,
) interpreter.TExecutedFunction {
	standardListenQuery := func(
		queryHandle *ttsapi.TQueryHandle,
		execInfo interpreter.TExecutionInfo,
	) interface{} {
		response, err := ttsApi.ListenQuery(queryHandle)
		if err != nil {
			return fmt.Errorf(
				"!error (line=%v, char=%v): %v",
				execInfo.Line,
				execInfo.CharPos,
				err.Error(),
			)
		}

		return []byte(*response)
	}

	return func(
		staticArgs interpreter.TFunctionArgumentsTable,
		_inputs interpreter.TFunctionInputChannelTable,
		internalGlobals interpreter.TGlobalVariablesTable,
		_externalGlobals interpreter.TGlobalVariablesTable,
		execInfo interpreter.TExecutionInfo,
		_interpreter *interpreter.Interpreter,
	) interface{} {
		fromVar, err := getFromVar(staticArgs, execInfo)
		if err != nil {
			return err
		}

		rawQueryHandle, hasQueryHandle := internalGlobals[fromVar]
		if !hasQueryHandle {
			return fmt.Errorf(
				"!error (line=%v, char=%v): query handle by name \"%v\" doesn't exist in internal variables",
				execInfo.Line,
				execInfo.CharPos,
				fromVar,
			)
		}

		queryHandle, isQueryHandleValid := rawQueryHandle.(*ttsapi.TQueryHandle)
		if !isQueryHandleValid {
			return fmt.Errorf(
				"!error (line=%v, char=%v): query handle by name \"%v\" is not valid as it's %v",
				execInfo.Line,
				execInfo.CharPos,
				fromVar,
				queryHandle,
			)
		}
		return standardListenQuery(queryHandle, execInfo)
	}
}
