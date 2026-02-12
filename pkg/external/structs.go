package external

type DalleImages struct {
	B64_JSON string `json:"b64_json"`
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

type OpenAIChatMessage struct {
	Role    string `json:"role"`
	Content any    `json:"content"`
}

type OpenAIChatContentPart struct {
	Type     string              `json:"type"`
	Text     string              `json:"text,omitempty"`
	ImageURL *OpenAIChatImageURL `json:"image_url,omitempty"`
}

type OpenAIChatImageURL struct {
	URL string `json:"url"`
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

type Response struct {
	ID               string         `json:"id"`
	Object           string         `json:"object"`
	CreatedAt        int64          `json:"created_at"`
	Status           string         `json:"status"`
	Background       bool           `json:"background"`
	Billing          Billing        `json:"billing"`
	Error            *Error         `json:"error"`
	Incomplete       any            `json:"incomplete_details"`
	Instructions     any            `json:"instructions"`
	MaxOutputTokens  *int           `json:"max_output_tokens"`
	MaxToolCalls     *int           `json:"max_tool_calls"`
	Model            string         `json:"model"`
	Output           []Output       `json:"output"`
	ParallelToolCall bool           `json:"parallel_tool_calls"`
	Reasoning        Reasoning      `json:"reasoning"`
	Usage            Usage          `json:"usage"`
	ServiceTier      string         `json:"service_tier"`
	Text             Text           `json:"text"`
	Tools            []Tool         `json:"tools"`
	User             *string        `json:"user"`
	Metadata         map[string]any `json:"metadata"`
}

type Error struct {
	Message string `json:"message"`
	Type    string `json:"type"`
	Param   string `json:"param"`
	Code    string `json:"code"`
}

type Billing struct {
	Payer string `json:"payer"`
}

type Output struct {
	ID            string          `json:"id"`
	Type          string          `json:"type"`
	Status        string          `json:"status"`
	Background    string          `json:"background,omitempty"`
	OutputFormat  string          `json:"output_format,omitempty"`
	Quality       string          `json:"quality,omitempty"`
	Result        string          `json:"result,omitempty"`
	RevisedPrompt string          `json:"revised_prompt,omitempty"`
	Size          string          `json:"size,omitempty"`
	Content       []OutputContent `json:"content,omitempty"`
	Role          string          `json:"role,omitempty"`
}

type OutputContent struct {
	Type        string `json:"type"`
	Text        string `json:"text"`
	Logprobs    []any  `json:"logprobs"`
	Annotations []any  `json:"annotations"`
}

type Reasoning struct {
	Effort  string `json:"effort"`
	Summary any    `json:"summary"`
}

type Usage struct {
	InputTokens  int `json:"input_tokens"`
	OutputTokens int `json:"output_tokens"`
	TotalTokens  int `json:"total_tokens"`
}

type Text struct {
	Format    TextFormat `json:"format"`
	Verbosity string     `json:"verbosity"`
}

type TextFormat struct {
	Type string `json:"type"`
}

type Tool struct {
	Type              string `json:"type"`
	Background        string `json:"background"`
	Moderation        string `json:"moderation"`
	N                 int    `json:"n"`
	OutputCompression int    `json:"output_compression"`
	OutputFormat      string `json:"output_format"`
	Quality           string `json:"quality"`
	Size              string `json:"size"`
}

type VideoJob struct {
	CompletedAt        int64     `json:"completed_at"`          // Unix timestamp for when the job completed
	CreatedAt          int64     `json:"created_at"`            // Unix timestamp for when the job was created
	Error              *JobError `json:"error,omitempty"`       // Error details if generation failed
	ExpiresAt          int64     `json:"expires_at"`            // Unix timestamp for when assets expire
	ID                 string    `json:"id"`                    // Unique identifier for the video job
	Model              string    `json:"model"`                 // Generation model used
	Object             string    `json:"object"`                // Object type (always "video")
	Progress           int       `json:"progress"`              // Completion percentage
	RemixedFromVideoID string    `json:"remixed_from_video_id"` // ID of source video if it's a remix
	Seconds            string    `json:"seconds"`               // Duration of generated clip (in seconds)
	Size               string    `json:"size"`                  // Video resolution
	Status             string    `json:"status"`                // Current lifecycle status
}

// If "error" is an object with details, define it:
type JobError struct {
	Message string `json:"message,omitempty"`
	Code    string `json:"code,omitempty"`
}
