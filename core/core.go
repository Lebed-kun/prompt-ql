package interpretercore

import (
	"fmt"
	"strings"
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
}

func New(
	execFnTable TExecutedFunctionTable,
	dataSwitchFn TDataSwitchFunction,
	defaultExternalVars TGlobalVariablesTable,
	defaultExternalVarsMeta TExternalGlobalsMetaTable,
) *Interpreter {
	execCtxStack := []*TExecutionStackFrame{
		makeRootStackFrame(),
	}

	return &Interpreter{
		// Settings
		execFnTable:   execFnTable,
		dataSwitchFn:  dataSwitchFn,
		defaultExternalGlobals: defaultExternalVars,
		defaultExternalGlobalsMeta: defaultExternalVarsMeta,

		// Internal interpreter state
		mode:          InterpreterModePlainText,
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
	}
}

// [BEGIN] Basic API

func (self *Interpreter) ExecutePartial(program string) *TInterpreterResult {
	res := self.executeImpl([]rune(program))
	self.resetPosition()
	return res
}

func (self *Interpreter) Execute(program string) *TInterpreterResult {
	res := self.executeImpl([]rune(program))
	if self.sessionClosed {
		self.resetImpl()
	} else {
		self.resetPosition()
	}
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

	embdRunes := []rune(embd)
	res := strings.Builder{}
	ptr := 0
	for ptr < len(embdRunes) {
		if embdRunes[ptr] == '%' && ptr < len(embdRunes) - 1 && isAlphaChar(embdRunes[ptr+1]) {
			ptr++

			begin := ptr
			for ptr < len(embdRunes) && isAlphaChar(embdRunes[ptr]) {
				ptr++
			}

			argName := string(embdRunes[begin:ptr])
			argVal, hasArg := args[argName]
			if !hasArg {
				argVal = fmt.Sprintf("%%%v", argName)
			}

			res.WriteString(argVal)
			continue
		}

		res.WriteRune(embdRunes[ptr])
		ptr++
	}

	return res.String(), nil
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
}

// [END] Misc control flow API
