package interpretercore

type TFunctionArgumentsTable map[string]interface{}

type TFunctionInputChannel []interface{}

type TFunctionInputChannelTable map[string]TFunctionInputChannel

type TExecutionStackFrame struct {
	State         int
	FnName        string
	ArgsTable     TFunctionArgumentsTable
	InputChannels TFunctionInputChannelTable
}

type TExecutionStack []*TExecutionStackFrame

type TExecutedFunction func(globalVars TGlobalVariablesTable, staticArgs TFunctionArgumentsTable, inputs TFunctionInputChannelTable) interface{}

type TExecutedFunctionTable map[string]TExecutedFunction

type TDataSwitchFunction func(topCtx *TExecutionStackFrame, rawData interface{})

type TGlobalVariablesTable map[string]interface{}

type TInterpreterResult struct {
	Result   interface{}
	Error    error
	Finished bool
}
