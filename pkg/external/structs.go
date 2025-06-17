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
	ID      string                `json:"id"`
	Object  string                `json:"object"`
	Model   string                `json:"model"`
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

type PerplexityRequest struct {
	Model                  string            `json:"model"`
	Message                PerplexityMessage `json:"message"`
	MaxTokens              int               `json:"max_tokens"`
	Temperature            float64           `json:"temperature"`
	TopP                   float64           `json:"top_p"`
	SearchDomainFilter     []string          `json:"search_domain_filter"`
	ReturnImages           bool              `json:"return_images"`
	ReturnRelatedQuestions bool              `json:"return_related_questions"`
	SearchRecencyFilter    string            `json:"search_recency_filter"`
	TopK                   int               `json:"top_k"`
	Stream                 bool              `json:"stream"`
	PresencePenalty        float64           `json:"presence_penalty"`
	FrequencyPenalty       float64           `json:"frequency_penalty"`
	ResponseFormat         string            `json:"response_format"`
}

type PerplexityMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type PerplexityResponse struct {
	ID            string                   `json:"id"`
	Model         string                   `json:"model"`
	Object        string                   `json:"object"`
	Created       int                      `json:"created"`
	Citations     []string                 `json:"citations"`
	SearchResults []PerplexitySearchResult `json:"search_results"`
	Choices       []PerplexityChoice       `json:"choices"`
	Usage         PerplexityUsage          `json:"usage"`
}

type PerplexityChoice struct {
	Index        int               `json:"index"`
	FinishReason string            `json:"finish_reason"`
	Message      PerplexityMessage `json:"message"`
	Delta        PerplexityMessage `json:"delta"`
}

type PerplexityUsage struct {
	PromptTokens      int    `json:"prompt_tokens"`
	CompletionTokens  int    `json:"completion_tokens"`
	TotalTokens       int    `json:"total_tokens"`
	SearchContextSize string `json:"search_context_size"`
}

type PerplexitySearchResult struct {
	Title string `json:"title"`
	URL   string `json:"url"`
	Date  string `json:"date"`
}
