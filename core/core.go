package interpretercore

import (
	"fmt"
	"strings"
)

type Interpreter struct {
	mode          int
	line          int
	charPos       int
	strPos        int
	globals       TGlobalVariablesTable
	execCtxStack  TExecutionStack
	isDirty       bool
	criticalError error
	execFnTable   TExecutedFunctionTable
	dataSwitchFn  TDataSwitchFunction
}

func makeRootStackFrame() *TExecutionStackFrame {
	inputChannels := make(TFunctionInputChannelTable)
	inputChannels["data"] = make(TFunctionInputChannel, 0)
	inputChannels["error"] = make(TFunctionInputChannel, 0)

	return &TExecutionStackFrame{
		State:         StackFrameStateExpectCmd,
		FnName:        "root",
		ArgsTable:     make(TFunctionArgumentsTable),
		InputChannels: inputChannels,
	}
}

func New(
	execFnTable TExecutedFunctionTable,
	dataSwitchFn TDataSwitchFunction,
) *Interpreter {
	execCtxStack := []*TExecutionStackFrame{
		makeRootStackFrame(),
	}

	return &Interpreter{
		mode:    InterpreterModePlainText,
		line:    0,
		charPos: 0,
		strPos:  0,
		globals: make(TGlobalVariablesTable),
		execCtxStack: execCtxStack,
		isDirty:       false,
		criticalError: nil,
		execFnTable:   execFnTable,
		dataSwitchFn:  dataSwitchFn,
	}
}

func (self *Interpreter) resetImpl() {
	self.mode = InterpreterModePlainText
	self.line = 0
	self.charPos = 0
	self.strPos = 0
	self.globals = make(TGlobalVariablesTable)
	self.execCtxStack = []*TExecutionStackFrame{
		makeRootStackFrame(),
	}
	self.isDirty = false
	self.criticalError = nil
}

func (self *Interpreter) getError(errorDetails string) error {
	return fmt.Errorf(
		"ERROR (line=%v, charpos=%v): %v",
		self.line,
		self.charPos,
		errorDetails,
	)
}

func (self *Interpreter) handlePlainText(program []rune) {
	plainText := strings.Builder{}
	escapeChar := false

	for self.strPos < len(program) && (program[self.strPos] != '{' || escapeChar) {
		if !escapeChar && program[self.strPos] == '\\' {
			escapeChar = true
		} else {
			escapeChar = false
			plainText.WriteRune(program[self.strPos])
		}
		self.strPos++

		if program[self.strPos-1] == '\n' {
			self.line++
			self.charPos = 0
		} else {
			self.charPos++
		}
	}

	topCtx := self.execCtxStack[len(self.execCtxStack)-1]
	self.dataSwitchFn(
		topCtx,
		fmt.Sprintf("!data %v", plainText.String()),
	)
}

func (self *Interpreter) resolveVariable(program []rune) (interface{}, error) {
	applyCnt := 1

	for self.strPos < len(program) && program[self.strPos] == '$' {
		applyCnt++
		self.strPos++
		self.charPos++
	}

	if self.strPos >= len(program) || !isAlphaNumChar(program[self.strPos]) {
		return nil, self.getError(
			"expected variable name in lating letters, digits or _ after $",
		)
	}

	begin := self.strPos
	for self.strPos < len(program) && isAlphaNumChar(program[self.strPos]) {
		self.strPos++
		self.charPos++
	}

	var varValue interface{}
	varName := string(program[begin:self.strPos])
	currApplyCnt := 0

	for currApplyCnt < applyCnt {
		var hasVar bool
		varValue, hasVar = self.globals[varName]
		if !hasVar {
			varValue = nil
			break
		}

		currApplyCnt++
		nextVarName, isNextVarNameStr := varValue.(string)
		if !isNextVarNameStr {
			break
		}
		varName = nextVarName
	}

	if varValue != nil && currApplyCnt != applyCnt {
		return nil, self.getError(
			"variable name is not string",
		)
	}

	return varValue, nil
}

func (self *Interpreter) resolveStrLiteral(program []rune) (string, error) {
	literal := strings.Builder{}
	escapeChar := false

	for self.strPos < len(program) && (program[self.strPos] != '"' || escapeChar) {
		if !escapeChar && program[self.strPos] == '\\' {
			escapeChar = true
		} else {
			escapeChar = false
			literal.WriteRune(program[self.strPos])
		}
		self.strPos++

		if program[self.strPos-1] == '\n' {
			self.line++
			self.charPos = 0
		} else {
			self.charPos++
		}
	}

	if self.strPos == len(program) && program[self.strPos-1] != '"' {
		return "", self.getError(
			"string literal must end with \"",
		)
	}

	self.strPos++
	return literal.String(), nil
}

func (self *Interpreter) skipWhitespaces(program []rune) {
	for self.strPos < len(program) && isWhitespace(program[self.strPos]) {
		if program[self.strPos] == '\n' {
			self.charPos = 0
			self.line++
		} else {
			self.charPos++
		}

		self.strPos++
	}
}

func (self *Interpreter) resolveName(program []rune) string {
	name := strings.Builder{}

	for self.strPos < len(program) && isAlphaNumChar(program[self.strPos]) {
		name.WriteRune(program[self.strPos])
		self.strPos++
		self.charPos++
	}

	return name.String()
}

func (self *Interpreter) handleCommand(program []rune) {
	var currLiteral interface{} = nil
	currArg := ""

	for self.strPos < len(program) && program[self.strPos] != '}' {
		if self.criticalError != nil {
			break
		}

		if isWhitespace(program[self.strPos]) {
			self.skipWhitespaces(program)
			continue
		}

		if program[self.strPos] == '~' {
			newCtxFrame := &TExecutionStackFrame{
				State:         StackFrameStateExpectCmd,
				FnName:        "",
				ArgsTable:     make(TFunctionArgumentsTable),
				InputChannels: make(TFunctionInputChannelTable, 0),
			}
			self.execCtxStack = append(self.execCtxStack, newCtxFrame)

			self.strPos++
			self.charPos++
			continue
		}

		if program[self.strPos] == '/' {
			self.strPos++
			self.charPos++

			topCtx := self.execCtxStack[len(self.execCtxStack)-1]
			topCtx.State = StackFrameStateIsClosing

			continue
		}

		if program[self.strPos] == '=' {
			topCtx := self.execCtxStack[len(self.execCtxStack)-1]
			if topCtx.State != StackFrameStateExpectArg {
				self.criticalError = self.getError(
					"expected argument before = ",
				)
				continue
			}

			currLiteralStr, isCurrLiteralStr := currLiteral.(string)
			if !isCurrLiteralStr {
				self.criticalError = self.getError(
					"argument name is not string",
				)
				continue
			}

			self.strPos++
			self.charPos++
			currArg = currLiteralStr
			topCtx.State = StackFrameStateExpectVal
			continue
		}

		if program[self.strPos] == '$' {
			self.strPos++
			self.charPos++

			var err error
			currLiteral, err = self.resolveVariable(program)

			if err != nil {
				self.criticalError = err
				continue
			}

			goto ctxFill
		}

		if program[self.strPos] == '"' {
			self.strPos++
			self.charPos++

			var err error
			currLiteral, err = self.resolveStrLiteral(program)

			if err != nil {
				self.criticalError = err
				continue
			}

			goto ctxFill
		}

		if isAlphaChar(program[self.strPos]) {
			currLiteral = self.resolveName(program)
			goto ctxFill
		}

		goto cmdParseError

	ctxFill:
		{
			topCtx := self.execCtxStack[len(self.execCtxStack)-1]
			switch topCtx.State {
			case StackFrameStateExpectCmd:
				currLiteralStr, isCurrLiteralStr := currLiteral.(string)
				if !isCurrLiteralStr {
					self.criticalError = self.getError(
						"command name is not string",
					)
				} else {
					topCtx.FnName = currLiteralStr
					topCtx.State = StackFrameStateExpectArg
				}
			case StackFrameStateExpectArg:
				if len(currArg) > 0 {
					topCtx.ArgsTable[currArg] = true
					currArg = ""
				} else {
					currLiteralStr, isCurrLiteralStr := currLiteral.(string)
					if !isCurrLiteralStr {
						self.criticalError = self.getError(
							"argument name is not string",
						)
					} else {
						currArg = currLiteralStr
					}
				}
			case StackFrameStateExpectVal:
				if len(currArg) == 0 {
					self.criticalError = self.getError(
						"argument is not provided",
					)
				} else {
					topCtx.ArgsTable[currArg] = currLiteral
					topCtx.State = StackFrameStateExpectArg
				}
			}

			continue
		}

	cmdParseError:
		{
			self.criticalError = self.getError(
				"unknown character",
			)
		}
	}
}

func (self *Interpreter) resolveTopCtx() {
	topCtx := self.execCtxStack[len(self.execCtxStack)-1]
	if len(self.execCtxStack) < 2 || topCtx.State != StackFrameStateIsClosing {
		return
	}

	self.execCtxStack = self.execCtxStack[:len(self.execCtxStack)-1]

	cmd, hasCmd := self.execFnTable[topCtx.FnName]
	if !hasCmd {
		self.criticalError = self.getError(
			fmt.Sprintf(
				"command with name \"%v\" doesn't exist in interpreter table",
				topCtx.FnName,
			),
		)
		return
	}

	if errChan, hasErrChan := topCtx.InputChannels["error"]; hasErrChan && len(errChan) > 0 {
		_, hasTopErrChan := self.execCtxStack[len(self.execCtxStack)-1].InputChannels["error"]

		if !hasTopErrChan {
			self.execCtxStack[len(self.execCtxStack)-1].InputChannels["error"] = make(TFunctionInputChannel, 0)
		}

		topErrChan := self.execCtxStack[len(self.execCtxStack)-1].InputChannels["error"]
		self.execCtxStack[len(self.execCtxStack)-1].InputChannels["error"] = append(
			topErrChan,
			topCtx.InputChannels["error"]...,
		)
	} else {
		result := cmd(
			self.globals,
			topCtx.ArgsTable,
			topCtx.InputChannels,
		)
		self.dataSwitchFn(
			self.execCtxStack[len(self.execCtxStack)-1],
			result,
		)
	}
}

func (self *Interpreter) executeImpl(program []rune) *TInterpreterResult {
	self.isDirty = true

	for self.strPos < len(program) {
		if self.criticalError != nil {
			break
		}

		switch program[self.strPos] {
		case '{':
			self.mode = InterpreterModeCommand
			self.strPos++
			self.charPos++
		case '}':
			self.resolveTopCtx()
			self.mode = InterpreterModePlainText
			self.strPos++
			self.charPos++
		}

		switch self.mode {
		case InterpreterModePlainText:
			self.handlePlainText(program)
		case InterpreterModeCommand:
			self.handleCommand(program)
		}
	}

	topCtx := self.execCtxStack[len(self.execCtxStack)-1]
	finished := self.criticalError != nil ||
		len(self.execCtxStack) == 1

	topErrChan, hasTopErrChan := topCtx.InputChannels["error"]
	if hasTopErrChan && len(topErrChan) > 0 {
		finished = true
	}

	return &TInterpreterResult{
		Result:   topCtx.InputChannels,
		Error:    self.criticalError,
		Finished: finished,
	}
}

func (self *Interpreter) ExecutePartial(program string) *TInterpreterResult {
	return self.executeImpl([]rune(program))
}

func (self *Interpreter) Execute(program string) *TInterpreterResult {
	res := self.executeImpl([]rune(program))
	self.resetImpl()
	return res
}

func (self *Interpreter) Reset() {
	self.resetImpl()
}
