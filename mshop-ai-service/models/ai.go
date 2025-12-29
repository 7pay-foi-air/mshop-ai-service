package models

type AIRequest struct {
	Prompt    string `json:"prompt"`
	MaxTokens int    `json:"max_tokens,omitempty"`
}

type AIResponse struct {
	Response string `json:"response"`
}
