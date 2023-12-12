package interpretercore

func (self *Interpreter) resetImpl() {
	self.resetPosition()
	self.mode = interpreterModePlainText
	self.internalGlobals = make(TGlobalVariablesTable)
	self.externalGlobals = initializeGlobals(self.defaultExternalGlobals)
	self.execCtxStack = []*TExecutionStackFrame{
		makeRootStackFrame(),
	}
	self.isDirty = false
	self.criticalError = nil
	self.embeddings = make(TEmbeddingsTable)
}

func (self *Interpreter) resetPosition() {
	self.line = 0
	self.charPos = 0
	self.strPos = 0
}