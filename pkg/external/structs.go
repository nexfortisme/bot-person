package external

type DalleImages struct {
	URL string `json:"url"`
}

type DalleResponse struct {
	Created int           `json:"created"`
	Data    []DalleImages `json:"data"`
}

type OpenAIDivinciChoice struct {
	Text   string `json:"text"`
	Index  int    `json:"index"`
	Reason string `json:"finish_reason"`
}

type OpenAIDivinciResponse struct {
	ID      string         `json:"id"`
	Object  string         `json:"object"`
	Model   string         `json:"model"`
	Choices []OpenAIDivinciChoice `json:"choices"`
}

type OpenAIGPTChoice struct {
	Index         int              `json:"index"`
	Message       OpenAIGPTMessage `json:"message"`
	Finish_Reason string           `json:"finish_reason"`
	LogProbs      string           `json:"logprobs"`
}

type OpenAIGPTMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type OpenAIGPTResponse struct {
	ID      string            `json:"id"`
	Object  string            `json:"object"`
	Created int64             `json:"created"`
	Model   string            `json:"model"`
	Usage   OpenAIGPTUsage    `json:"usage"`
	Choices []OpenAIGPTChoice `json:"choices"`
}

type OpenAIGPTUsage struct {
	Prompt_Tokens     int `json:"prompt_tokens"`
	Completion_Tokens int `json:"completion_tokens"`
	Total_Tokens      int `json:"total_tokens"`
}