package interpretercore

type InterpreterApi interface {
	Execute(program string) *TInterpreterResult
	UnsafeExecute(program string) *TInterpreterResult
	Reset()
	IsDirty() bool
	IsSessionClosed() bool
}
