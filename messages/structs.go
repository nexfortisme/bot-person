package messages

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

type SDResponse struct {
	Status string `json:"status"`
	Request interface{} `json:"request"`
	Output []SDOutput `json:"output"`
}

type SDStep struct {
	step int `json:"step"`
	total_steps int `json:"total_steps"`
}

type SDOutput struct {
	Data string `json:"data"`
	Seed int `json:"seed"`
	Path_abs interface{} `json:"path_abs"`
}
