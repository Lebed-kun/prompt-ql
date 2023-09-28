package interpretercore

type Interpreter struct {
	// Settings
	execFnTable   TExecutedFunctionTable
	dataSwitchFn  TDataSwitchFunction
	defaultExternalGlobals TGlobalVariablesTable

	// Internal interpreter state
	mode          int
	line          int
	charPos       int
	strPos        int
	isDirty       bool
	criticalError error
	execCtxStack  TExecutionStack
	
	// Current globals
	internalGlobals       TGlobalVariablesTable
	externalGlobals TGlobalVariablesTable
}

func New(
	execFnTable TExecutedFunctionTable,
	dataSwitchFn TDataSwitchFunction,
	defaultExternalVars TGlobalVariablesTable,
) *Interpreter {
	execCtxStack := []*TExecutionStackFrame{
		makeRootStackFrame(),
	}

	return &Interpreter{
		// Settings
		execFnTable:   execFnTable,
		dataSwitchFn:  dataSwitchFn,
		defaultExternalGlobals: defaultExternalVars,

		// Internal interpreter state
		mode:          InterpreterModePlainText,
		line:          0,
		charPos:       0,
		strPos:        0,
		execCtxStack:  execCtxStack,
		isDirty:       false,
		criticalError: nil,

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
	self.resetImpl()
	return res
}

func (self *Interpreter) Reset() {
	self.resetImpl()
}

func (self *Interpreter) IsDirty() bool {
	return self.isDirty
}

func (self *Interpreter) SetExternalGlobals(globals TGlobalVariablesTable) {
	self.defaultExternalGlobals = globals
	self.externalGlobals = initializeExternalGlobals(self.defaultExternalGlobals)
}

func (self *Interpreter) SetExternalGlobalVar(name string, val interface{}) {
	if self.defaultExternalGlobals == nil {
		self.defaultExternalGlobals = make(TGlobalVariablesTable)
		self.externalGlobals = initializeExternalGlobals(self.defaultExternalGlobals)
	}
	
	self.defaultExternalGlobals[name] = val
	self.externalGlobals[name] = val
}
