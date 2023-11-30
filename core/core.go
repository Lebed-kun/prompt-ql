package interpretercore

import (
	"fmt"
)

type Interpreter struct {
	// Settings
	execFnTable   TExecutedFunctionTable
	dataSwitchFn  TDataSwitchFunction
	defaultExternalGlobals TGlobalVariablesTable
	defaultExternalGlobalsMeta TExternalGlobalsMetaTable

	// Internal interpreter state
	mode          int
	line          int
	charPos       int
	strPos        int
	isDirty       bool
	criticalError error
	execCtxStack  TExecutionStack
	sessionClosed bool
	
	// Current globals
	internalGlobals       TGlobalVariablesTable
	externalGlobals TGlobalVariablesTable

	// Current embeddings
	embeddings TEmbeddingsTable
	embeddingsMeta TEmbeddingMetaInfoTable

	// Current restrictions
	restrictedCmds TRestrictedCommands
}

func New(
	execFnTable TExecutedFunctionTable,
	dataSwitchFn TDataSwitchFunction,
	defaultExternalVars TGlobalVariablesTable,
	defaultExternalVarsMeta TExternalGlobalsMetaTable,
	restrictedCommands TRestrictedCommands,
) *Interpreter {
	execCtxStack := []*TExecutionStackFrame{
		makeRootStackFrame(),
	}

	if restrictedCommands == nil {
		restrictedCommands = make(TRestrictedCommands)
	}

	return &Interpreter{
		// Settings
		execFnTable:   execFnTable,
		dataSwitchFn:  dataSwitchFn,
		defaultExternalGlobals: defaultExternalVars,
		defaultExternalGlobalsMeta: defaultExternalVarsMeta,

		// Internal interpreter state
		mode:          interpreterModePlainText,
		line:          0,
		charPos:       0,
		strPos:        0,
		execCtxStack:  execCtxStack,
		isDirty:       false,
		criticalError: nil,
		sessionClosed: true,

		// Current globals
		internalGlobals:       make(TGlobalVariablesTable),
		externalGlobals: initializeExternalGlobals(defaultExternalVars),

		// Current embeddings
		embeddings: make(TEmbeddingsTable),
		embeddingsMeta: make(TEmbeddingMetaInfoTable),

		// Current restrictions
		restrictedCmds: restrictedCommands,
	}
}

// [BEGIN] Basic API

func (self *Interpreter) Execute(program string) *TInterpreterResult {
	res := self.executeFullImpl([]rune(program), true)
	return res
}

func (self *Interpreter) UnsafeExecute(program string) *TInterpreterResult {
	res := self.executeFullImpl([]rune(program), false)
	return res
}

func (self *Interpreter) Reset() {
	self.resetImpl()
}

func (self *Interpreter) IsDirty() bool {
	return self.isDirty
}

// [END] Basic API

// [BEGIN] Globals API

func (self *Interpreter) SetExternalGlobals(globals TGlobalVariablesTable, globalsMeta TExternalGlobalsMetaTable) {
	self.defaultExternalGlobals = globals
	self.defaultExternalGlobalsMeta = globalsMeta
	self.externalGlobals = initializeExternalGlobals(self.defaultExternalGlobals)
}

func (self *Interpreter) SetExternalGlobalVar(name string, val interface{}, description string) {
	if self.defaultExternalGlobals == nil {
		self.defaultExternalGlobals = make(TGlobalVariablesTable)
		self.externalGlobals = initializeExternalGlobals(self.defaultExternalGlobals)
	}
	if self.defaultExternalGlobalsMeta == nil && len(description) > 0 {
		self.defaultExternalGlobalsMeta = make(TExternalGlobalsMetaTable)
	}
	
	self.defaultExternalGlobals[name] = val
	self.externalGlobals[name] = val
	if self.defaultExternalGlobalsMeta != nil && len(description) > 0 {
		self.defaultExternalGlobalsMeta[name] = &TExternalGlobalMetaInfo{
			Description: description,
		}
	}
}

func (self *Interpreter) GetExternalGlobalsList() map[string]string {
	res := make(map[string]string, 0)
	
	if self.defaultExternalGlobalsMeta == nil {
		for k := range self.defaultExternalGlobals {
			res[k] = ""
		}
	} else {
		for k, v := range self.defaultExternalGlobalsMeta {
			res[k] = v.Description
		}
	}

	return res
}

// [END] Globals API

// [BEGIN] Sessions API

func (self *Interpreter) OpenSession() {
	self.sessionClosed = false
}

func (self *Interpreter) CloseSession() {
	self.sessionClosed = true
}

func (self *Interpreter) IsSessionClosed() bool {
	return self.sessionClosed
}

// [END] Sessions API

// [BEGIN] Embeddings API

func (self *Interpreter) GetEmbeddingsList() map[string]string {
	res := make(map[string]string, 0)
	
	for k, v := range self.embeddingsMeta {
		res[k] = v.Description
	}

	return res
}

func (self *Interpreter) RegisterEmbedding(name string, code string, description string) {
	self.embeddings[name] = code
	if len(description) > 0 {
		self.embeddingsMeta[name] = &TEmbeddingMetaInfo{
			Description: description,
		}
	}
}

func (self *Interpreter) ExpandEmbedding(name string, args TEmbeddingArgsTable) (string, error) {
	embd, hasEmbd := self.embeddings[name]
	if !hasEmbd {
		return "", self.getError(
			fmt.Sprintf("embedding named \"%v\" is not defined", name),
		)
	}

	return self.expandImpl(embd, args), nil
}

func (self *Interpreter) ExpandInlineEmbedding(embedding string, args TEmbeddingArgsTable) string {
	return self.expandImpl(embedding, args)
}

// [END] Embeddings API

// [BEGIN] Misc control flow API

func (self *Interpreter) ControlFlowClearInternalVars() {
	self.internalGlobals = make(TGlobalVariablesTable)
}

func (self *Interpreter) ControlFlowClearStack() {
	self.execCtxStack = []*TExecutionStackFrame{
		makeRootStackFrame(),
	}
	self.criticalError = nil
}

// [END] Misc control flow API
