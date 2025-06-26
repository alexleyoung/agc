package types

import "google.golang.org/genai"

type ChatRequestBody struct {
	Prompt  string           `json:"prompt"`
	Model   string           `json:"model,omitempty"`
	History []*genai.Content `json:"history,omitempty"`
}
