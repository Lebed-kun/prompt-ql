package interpretercore

type TFunctionArgumentsTable map[string]interface{}
type TFunctionInputChannelTable map[string]TFunctionInputChannel

type TExecutionStackFrame struct {
	State         int
	FnName        string
	ArgsTable     TFunctionArgumentsTable
	InputChannels TFunctionInputChannelTable
}

type TExecutionStack []*TExecutionStackFrame

type TExecutedFunction func(
	staticArgs TFunctionArgumentsTable,
	inputs TFunctionInputChannelTable,
	internalGlobals TGlobalVariablesTable,
	externalGlobals TGlobalVariablesTable,
	execInfo TExecutionInfo,
	interpreter *Interpreter,
) interface{}

type TExecutedFunctionTable map[string]TExecutedFunction

type TDataSwitchFunction func(topCtx *TExecutionStackFrame, rawData interface{})

type TGlobalVariablesTable map[string]interface{}

type TExternalGlobalsMetaTable map[string]*TExternalGlobalMetaInfo

type TExternalGlobalMetaInfo struct {
	Description string
}

type TExecutionInfo struct {
	Line int
	CharPos int
	StrPos int
}

type TEmbeddingsTable map[string]string

type TEmbeddingArgsTable map[string]string

type TEmbeddingMetaInfo struct {
	Description string
}

type TEmbeddingMetaInfoTable map[string]*TEmbeddingMetaInfo

type TRestrictedCommands map[string]bool

type TCommandMetaInfo struct {
	IsErrorTolerant bool
}

type TCommandMetaInfoTable map[string]*TCommandMetaInfo
