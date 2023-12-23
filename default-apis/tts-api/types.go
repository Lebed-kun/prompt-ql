package ttsapi

type TTtsApiResponse []byte

type TTtsApiRequest struct {
	Model string `json:"model"`
	Input string `json:"input"`
	Voice string `json:"voice"`
}

type TQueryHandle struct {
	ResultChan chan *TTtsApiResponse
	ErrChan    chan error
}
