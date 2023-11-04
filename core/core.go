package interpretercore

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
	}
}

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

func (self *Interpreter) OpenSession() {
	self.sessionClosed = false
}

func (self *Interpreter) CloseSession() {
	self.sessionClosed = true
}

func (self *Interpreter) IsSessionClosed() bool {
	return self.sessionClosed
}
