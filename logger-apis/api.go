package loggerapis

type LoggersMainApi interface {
	RegisterLogger(
		name string,
		logger TLoggerFunc,
	)
	UnregisterLogger(name string)
}
