package interpretercore

const (
	stackFrameStateExpectCmd = iota
	stackFrameStateExpectArg
	stackFrameStateExpectVal
	stackFrameStateIsClosing
	stackFrameStateFullfilled
	stackFrameStateExpectCmdAfterFullfill
)

const (
	interpreterModePlainText = iota
	interpreterModeCommand
	interpreterModeCodeLiteral
	interpreterModeCodeComment
)
