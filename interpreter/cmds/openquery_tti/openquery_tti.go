package openquerytticmd

import (
	"fmt"
	interpreter "gitlab.com/jbyte777/prompt-ql/v5/core"
	ttiapi "gitlab.com/jbyte777/prompt-ql/v5/default-apis/tti-api"
	utils "gitlab.com/jbyte777/prompt-ql/v5/interpreter/utils"
)

func MakeOpenQueryTtiCmd(
	ttiApi *ttiapi.TtiApi,
) interpreter.TExecutedFunction {
	standardOpenQuery := func(
		model string,
		width uint,
		height uint,
		responseFormat string,
		inputs interpreter.TFunctionInputChannelTable,
		execInfo interpreter.TExecutionInfo,
		isSync bool,
	) (interface{}, error) {
		prompt, err := inputs["data"].MergeIntoString()
		if err != nil {
			return nil, err
		}

		if isSync {
			handle, err := ttiApi.OpenQuery(
				model,
				prompt,
				width,
				height,
				responseFormat,
			)
			if err != nil {
				return nil, err
			}
			response, err := ttiApi.ListenQuery(handle)
			if err != nil {
				return nil, err
			}

			blob, err := utils.ReadTtiBlob(response)
			if err != nil {
				return nil, err
			}

			return blob, nil
		} else {
			return ttiApi.OpenQuery(
				model,
				prompt,
				width,
				height,
				responseFormat,
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

		width, err := getWidth(staticArgs, execInfo)
		if err != nil {
			return err
		}

		height, err := getHeight(staticArgs, execInfo)
		if err != nil {
			return err
		}

		responseFormat, err := getResponseFormat(staticArgs, execInfo)
		if err != nil {
			return err
		}

		queryHandleOrResponse, err := standardOpenQuery(
			model,
			width,
			height,
			responseFormat,
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
