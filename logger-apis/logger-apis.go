package loggerapis

import (
	"fmt"

	interpreter "gitlab.com/jbyte777/prompt-ql/v5/core"
)

type LoggerApis struct {
	loggers	TLoggerFuncTable
}

func New() *LoggerApis {
	return &LoggerApis{
		loggers: TLoggerFuncTable{
			"default": defaultLogger,
		},
	}
}

func (self *LoggerApis) RegisterLogger(
	name string,
	logger TLoggerFunc,
) {
	self.loggers[name] = logger
}

func (self *LoggerApis) UnregisterLogger(name string) {
	delete(self.loggers, name)
}

func (self *LoggerApis) Log(
	loggerName string,
	execInfo interpreter.TExecutionInfo,
	inputs interpreter.TFunctionInputChannelTable,
	internalGlobals interpreter.TGlobalVariablesTable,
	externalGlobals interpreter.TGlobalVariablesTable,
) error {
	logger, hasLogger := self.loggers[loggerName]
	if !hasLogger {
		return fmt.Errorf(
			"!error (line=%v, char=%v): custom logger named \"%v\" doesn't exist",
			execInfo.Line,
			execInfo.CharPos,
			loggerName,
		)
	}

	return logger(
		execInfo,
		inputs,
		internalGlobals,
		externalGlobals,
	)
}
