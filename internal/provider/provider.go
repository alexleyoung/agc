package provider

import "context"

type Provider interface {
	NewClient() Client
}

type Client interface {
	NewChat(ctx context.Context, model string, config *GenerationConfig, history []*Content) (Chat, error)
}

type Chat interface {
	SendMessage(ctx context.Context, prompt string) (*ContentResponse, error)
}

// GenerationConfig is a struct that contains configuration options for the generation process.
// e.g. Temperature, System Instructions, tools, etc.
type GenerationConfig struct{}

type Tool struct{}

// ContentResponse is a struct that contains the response from the model.
type ContentResponse struct{}

type Content struct{}
