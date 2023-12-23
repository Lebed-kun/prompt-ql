package chatapi

type TMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type TGptApiResponseChoice struct {
	Index        int      `json:"index"`
	Message      TMessage `json:"message"`
	FinishReason string   `json:"finish_reason"`
}

type TGptApiResponse struct {
	Choices []TGptApiResponseChoice `json:"choices"`
}

type TGptApiRequest struct {
	Model       string     `json:"model"`
	Messages    []TMessage `json:"messages"`
	Temperature float64    `json:"temperature"`
	N           int        `json:"n"`
}

type TQueryHandle struct {
	ResultChan chan *TGptApiResponse
	ErrChan    chan error
}
