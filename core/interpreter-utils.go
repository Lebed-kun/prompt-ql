package interpretercore

func makeRootStackFrame() *TExecutionStackFrame {
	inputChannels := make(TFunctionInputChannelTable)
	inputChannels["data"] = make(TFunctionInputChannel, 0)
	inputChannels["error"] = make(TFunctionInputChannel, 0)

	return &TExecutionStackFrame{
		State:         StackFrameStateExpectCmd,
		FnName:        "root",
		ArgsTable:     make(TFunctionArgumentsTable),
		InputChannels: inputChannels,
	}
}

func initializeExternalGlobals(defaultGlobals TGlobalVariablesTable) TGlobalVariablesTable {
	if defaultGlobals == nil {
		return nil
	}
	
	res := make(TGlobalVariablesTable)

	for k, v := range defaultGlobals {
		res[k] = v
	}

	return res
}
