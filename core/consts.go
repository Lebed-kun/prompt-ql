package interpretercore

const (
	StackFrameStateExpectCmd = iota
	StackFrameStateExpectArg
	StackFrameStateExpectVal
	StackFrameStateIsClosing
	StackFrameStateFullfilled
	StackFrameStateExpectCmdAfterFullfill
)

const (
	InterpreterModePlainText = iota
	InterpreterModeCommand
)
