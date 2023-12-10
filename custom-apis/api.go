package customapis

type CustomModelsMainApi interface {
	RegisterModelApi(
		name string,
		doQuery TDoQueryFunc,
		description string,
	)
	UnregisterModelApi(name string)
}
