package ttiapi

type TtiApiResponseEntry struct {
	B64Json string `json:"b64_json"`
	Url     string `json:"url"`
}

type TTtiApiResponse struct {
	Data []TtiApiResponseEntry `json:"data"`
}

type TTtiApiRequest struct {
	Model          string `json:"model"`
	Prompt         string `json:"prompt"`
	N              int    `json:"n"`
	Size           string `json:"size"`
	ResponseFormat string `json:"response_format"`
}

type TQueryHandle struct {
	ResultChan chan *TTtiApiResponse
	ErrChan    chan error
}
