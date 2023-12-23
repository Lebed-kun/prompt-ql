package openqueryttscmd

import (
	"fmt"
	ttsapi "gitlab.com/jbyte777/prompt-ql/v5/default-apis/tts-api"
	interpreter "gitlab.com/jbyte777/prompt-ql/v5/core"
)

func MakeOpenQueryTtsCmd(
	ttsApi *ttsapi.TtsApi,
) interpreter.TExecutedFunction {
	standardOpenQuery := func(
		model string,
		voice string,
		inputs interpreter.TFunctionInputChannelTable,
		execInfo interpreter.TExecutionInfo,
		isSync bool,
	) (interface{}, error) {
		prompt, err := inputs["data"].MergeIntoString()
		if err != nil {
			return nil, err
		}

		if isSync {
			handle, err := ttsApi.OpenQuery(
				model,
				prompt,
				voice,
			)
			if err != nil {
				return nil, err
			}
			response, err := ttsApi.ListenQuery(handle)
			if err != nil {
				return nil, err
			}
			return []byte(*response), nil
		} else {
			return ttsApi.OpenQuery(
				model,
				prompt,
				voice,
			)
		}
	}

	return func(
		staticArgs interpreter.TFunctionArgumentsTable,
		inputs interpreter.TFunctionInputChannelTable,
		internalGlobals interpreter.TGlobalVariablesTable,
		_externalGlobals interpreter.TGlobalVariablesTable,
		execInfo interpreter.TExecutionInfo,
		_interpreter *interpreter.Interpreter,
	) interface{} {
		syncFlag := getSyncFlag(staticArgs)

		toVar, err := getToVar(staticArgs, execInfo, syncFlag)
		if err != nil {
			return err
		}
	
		model, err := getModel(staticArgs, execInfo)
		if err != nil {
			return err
		}
		
		voice, err := getVoice(staticArgs, execInfo)
		if err != nil {
			return err
		}

		queryHandleOrResponse, err := standardOpenQuery(
			model,
			voice,
			inputs,
			execInfo,
			syncFlag,
		)

		if err != nil {
			return fmt.Errorf(
				"!error (line=%v, char=%v): %v",
				execInfo.Line,
				execInfo.CharPos,
				err.Error(),
			)
		}

		if !syncFlag {
			internalGlobals[toVar] = queryHandleOrResponse
		} else {
			return queryHandleOrResponse
		}

		return nil
	}
}
