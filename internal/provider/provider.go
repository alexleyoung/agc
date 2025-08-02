package provider

import "context"

type Client interface {
	GenerateContent(ctx context.Context, model string, contents []*Content, config *GenerationConfig)
}

type GenerationConfig struct{}

type Tool struct{}

type ContentResponse struct{}

type Content struct{}
