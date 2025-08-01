package provider

type Client interface {
	GenerateContent()
}

type GenerationConfig struct{}

type Tool struct{}

type ContentResponse struct{}
