package api

var supportedOpenAiModels map[string]bool = map[string]bool{
	"gpt-4": true,
	"gpt-4-0613": true,
	"gpt-4-32k": true,
	"gpt-4-32k-0613": true,
	"gpt-3.5-turbo": true,
	"gpt-3.5-turbo-0613": true,
	"gpt-3.5-turbo-16k": true,
	"gpt-3.5-turbo-16k-0613": true,
}
