package external

type OpenAIGPTChoice struct {
	Index         int              `json:"index"`
	Message       OpenAIGPTMessage `json:"message"`
	Finish_Reason string           `json:"finish_reason"`
	LogProbs      string           `json:"logprobs"`
}
