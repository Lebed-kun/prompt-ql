package interpreter

const (
	StackFrameStateExpectCmd = iota
	StackFrameStateExpectArg
	StackFrameStateExpectVal
	StackFrameStateIsClosing
)

const (
	InterpreterModePlainText = iota
	InterpreterModeCommand
)
