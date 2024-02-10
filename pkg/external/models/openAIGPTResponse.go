package external

type OpenAIGPTResponse struct {
	ID      string            `json:"id"`
	Object  string            `json:"object"`
	Created int64             `json:"created"`
	Model   string            `json:"model"`
	Usage   OpenAIGPTUsage    `json:"usage"`
	Choices []OpenAIGPTChoice `json:"choices"`
}