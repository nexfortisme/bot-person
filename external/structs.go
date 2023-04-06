package external

// The outer structure of the response from OpenAI
type OpenAIGPTResponse struct {
	ID      string            `json:"id"`
	Object  string            `json:"object"`
	Created int64             `json:"created"`
	Model   string            `json:"model"`
	Usage   OpenAIGPTUsage    `json:"usage"`
	Choices []OpenAIGPTChoice `json:"choices"`
}

// The inter structure of the response from OpenAI, this
// contains zero or more completions based on the provided
// prompt
type OpenAIGPTChoice struct {
	Index         int              `json:"index"`
	Message       OpenAIGPTMessage `json:"message"`
	Finish_Reason string           `json:"finish_reason"`
}

type OpenAIGPTUsage struct {
	Prompt_Tokens     int `json:"prompt_tokens"`
	Completion_Tokens int `json:"completion_tokens"`
	Total_Tokens      int `json:"total_tokens"`
}

type OpenAIGPTMessage struct {
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

// The outer structure of the response from OpenAI
type OpenAIResponse struct {
	ID      string         `json:"id"`
	Object  string         `json:"object"`
	Model   string         `json:"model"`
	Choices []OpenAIChoice `json:"choices"`
}

// The inter structure of the response from OpenAI, this
// contains zero or more completions based on the provided
// prompt
type OpenAIChoice struct {
	Text   string `json:"text"`
	Index  int    `json:"index"`
	Reason string `json:"finish_reason"`
}
