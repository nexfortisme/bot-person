package external

// The outer structure of the response from OpenAI
type OpenAIResponse struct {
	ID     string `json:"id"`
	Object string `json:"object"`
	Created string `json:"created"`
	Model string `json:"model"`
	Usage OpenAIUsage `json:"usage"`
	// Model   string         `json:"model"`
	Choices []OpenAIChoice `json:"choices"`
}

// The inter structure of the response from OpenAI, this
// contains zero or more completions based on the provided
// prompt
type OpenAIChoice struct {
	Index         string        `json:"index"`
	Message       OpenAIMessage `json:"message"`
	Finish_Reason string        `json:"finish_reason"`
	// Text   string `json:"text"`
	// Index  int    `json:"index"`
	// Reason string `json:"finish_reason"`
}

type OpenAIUsage struct {
	PromtTokens int `json:"prompt_tokens"`
	CompletionTokens int `json:"completion_tokens"`
	TotalTokens int `json:"total_tokens"`
}

type OpenAIMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type DalleResponse struct {
	Created int           `json:"created"`
	Data    []DalleImages `json:"data"`
}

type DalleImages struct {
	URL string `json:"url"`
}
