package interpretercore

func (self *Interpreter) resetImpl() {
	self.resetPosition()
	self.mode = InterpreterModePlainText
	self.internalGlobals = make(TGlobalVariablesTable)
	self.externalGlobals = initializeExternalGlobals(self.defaultExternalGlobals)
	self.execCtxStack = []*TExecutionStackFrame{
		makeRootStackFrame(),
	}
	self.isDirty = false
	self.criticalError = nil
}

func (self *Interpreter) resetPosition() {
	self.line = 0
	self.charPos = 0
	self.strPos = 0
}