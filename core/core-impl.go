package interpretercore

import (
	"fmt"
	"strings"
)

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

	rawText := plainText.String()
	topCtx := self.execCtxStack[len(self.execCtxStack)-1]
	self.dataSwitchFn(
		topCtx,
		fmt.Sprintf("!data %v", rawText),
	)
}

func (self *Interpreter) resolveReference(program []rune) (interface{}, error) {
	applyCnt := 1

	for self.strPos < len(program) && program[self.strPos] == '$' {
		applyCnt++
		self.strPos++
		self.charPos++
	}

	if self.strPos >= len(program) || (!isAlphaNumChar(program[self.strPos]) && program[self.strPos] != '@') {
		return nil, self.getError(
			"expected reference name in lating letters, digits, _ or @ after $",
		)
	}

	isReferenceExternal := false
	if program[self.strPos] == '@' {
		isReferenceExternal = true
		self.strPos++
		self.charPos++
	}

	if self.strPos >= len(program) || !isAlphaNumChar(program[self.strPos]) {
		return nil, self.getError(
			"expected reference name in lating letters, digits, or _ after $[@]",
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

	varTable := self.internalGlobals
	if isReferenceExternal {
		varTable = self.externalGlobals
	}

	for currApplyCnt < applyCnt {
		var hasVar bool
		varValue, hasVar = varTable[varName]
		if !hasVar {
			if isReferenceExternal {
				return nil, self.getError(
					fmt.Sprintf(
						"external variable with name \"%v\" is not defined",
						varName,
					),
				)
			}

			varValue = nil
			break
		}

		currApplyCnt++
		nextVarName, isNextVarNameStr := varValue.(string)
		if !isNextVarNameStr || len(nextVarName) == 0 {
			break
		}
		varName = nextVarName
	}

	if varValue != nil && currApplyCnt != applyCnt {
		return nil, self.getError(
			"reference name is a not string",
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

	if self.strPos < len(program) && program[self.strPos] == '@' {
		name.WriteRune(program[self.strPos])
		self.strPos++
		self.charPos++
	}

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

		topCtx := self.execCtxStack[len(self.execCtxStack)-1]
		if len(currArg) > 0 {
			topCtx.ArgsTable[currArg] = true
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
			currLiteral = nil
			currArg = ""
			continue
		}

		if program[self.strPos] == '/' {
			self.strPos++
			self.charPos++

			if topCtx.State == StackFrameStateFullfilled {
				topCtx.State = StackFrameStateExpectCmdAfterFullfill
			} else {
				topCtx.State = StackFrameStateIsClosing
			}

			currLiteral = nil
			currArg = ""
			continue
		}

		if program[self.strPos] == '=' {
			if topCtx.State != StackFrameStateExpectArg || len(currArg) == 0 {
				self.criticalError = self.getError(
					"expected argument before = ",
				)
				continue
			}

			self.strPos++
			self.charPos++
			currLiteral = nil
			topCtx.State = StackFrameStateExpectVal
			continue
		}

		if program[self.strPos] == '$' {
			self.strPos++
			self.charPos++

			var err error
			currLiteral, err = self.resolveReference(program)

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

		if isAlphaChar(program[self.strPos]) || program[self.strPos] == '@' {
			currLiteral = self.resolveName(program)
			goto ctxFill
		}

		goto cmdParseError

	ctxFill:
		{
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
					currArg = ""
				}

				currLiteralStr, isCurrLiteralStr := currLiteral.(string)
				if !isCurrLiteralStr {
					self.criticalError = self.getError(
						"argument name is not string",
					)
				} else {
					currArg = currLiteralStr
				}
			case StackFrameStateExpectVal:
				if len(currArg) == 0 {
					self.criticalError = self.getError(
						"argument is not provided",
					)
				} else {
					topCtx.ArgsTable[currArg] = currLiteral
					topCtx.State = StackFrameStateExpectArg
					currArg = ""
				}
			case StackFrameStateExpectCmdAfterFullfill:
				currLiteralStr, isCurrLiteralStr := currLiteral.(string)
				if !isCurrLiteralStr {
					self.criticalError = self.getError(
						"command name is not string",
					)
				} else if currLiteralStr != topCtx.FnName {
					self.criticalError = self.getError(
						fmt.Sprintf(
							"command \"%v\" does not match the command \"%v\" on the closest context",
							currLiteralStr,
							topCtx.FnName,
						),
					)
				} else {
					topCtx.State = StackFrameStateIsClosing
				}
			}

			continue
		}

	cmdParseError:
		{
			self.criticalError = self.getError(
				fmt.Sprintf(
					"unknown character %v",
					program[self.strPos],
				),
			)
		}
	}

	if self.criticalError != nil {
		return
	}

	topCtx := self.execCtxStack[len(self.execCtxStack)-1]
	if topCtx.State == StackFrameStateExpectCmdAfterFullfill {
		self.criticalError = self.getError(
			fmt.Sprintf(
				"closing command is empty while command on the closest context is \"%v\"",
				topCtx.FnName,
			),
		)
		return
	}
	
	if topCtx.State != StackFrameStateIsClosing {
		topCtx.State = StackFrameStateFullfilled
	}
	if len(currArg) > 0 {
		topCtx.ArgsTable[currArg] = true
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
			topCtx.ArgsTable,
			topCtx.InputChannels,
			self.internalGlobals,
			self.externalGlobals,
			TExecutionInfo{
				StrPos:  self.strPos,
				CharPos: self.charPos,
				Line:    self.line,
			},
			self,
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
	complete := self.criticalError != nil ||
		len(self.execCtxStack) == 1

	topErrChan, hasTopErrChan := topCtx.InputChannels["error"]
	if hasTopErrChan && len(topErrChan) > 0 {
		complete = true
	}

	return &TInterpreterResult{
		Result:   topCtx.InputChannels,
		Error:    self.criticalError,
		Complete: complete,
	}
}
