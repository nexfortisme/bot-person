package external

type OpenAIDivinciResponse struct {
	ID      string         `json:"id"`
	Object  string         `json:"object"`
	Model   string         `json:"model"`
	Choices []OpenAIDivinciChoice `json:"choices"`
}