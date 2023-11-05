package interpretercore

import (
	"fmt"
	"strings"
)

func (self *Interpreter) handlePlainText(program []rune) {
	plainText := strings.Builder{}
	escapeChar := false

	// end of program text
	for self.strPos < len(program) &&
		// start of command
		(escapeChar || program[self.strPos] != '{') &&
		// start of code literal
		(escapeChar || self.strPos >= len(program) - 1 || !(program[self.strPos] == '<' && program[self.strPos+1] == '%')) {
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

func (self *Interpreter) handleCodeLiteral(program []rune) {
	plainText := strings.Builder{}
	escapeChar := false

	// end of program text
	for self.strPos < len(program)-1 &&
		// end of code literal
		(escapeChar || !(program[self.strPos] == '%' && program[self.strPos+1] == '>')) {
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

	if self.strPos == len(program)-1 {
		self.criticalError = self.getError(
			"expected %> in the end of code literal",
		)
	} else {
		rawText := plainText.String()
		topCtx := self.execCtxStack[len(self.execCtxStack)-1]
		self.dataSwitchFn(
			topCtx,
			fmt.Sprintf("!data %v", rawText),
		)
	}
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