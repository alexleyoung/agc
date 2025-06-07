package ai

import (
	"context"
	"log"
	"os"

	"google.golang.org/genai"
)

var FunctionDeclarations = []*genai.FunctionDeclaration{{
	Name:        "create_event",
	Description: "Creates a new event in the user's calendar.",
	Parameters: &genai.Schema{
		Type: "object",
		Properties: map[string]*genai.Schema{
			"summary":     {Type: "string", Description: "The title of the event."},
			"description": {Type: "string", Description: "The description of the event."},
			"start":       {Type: "string", Description: "The time, as a combined date-time value (formatted according to RFC3339) with a timezone offset"},
			"end":         {Type: "string", Description: "The time, as a combined date-time value (formatted according to RFC3339) with a timezone offset"},
		},
	},
}}

func Chat(ctx context.Context, model string, prompt string) *genai.GenerateContentResponse {
	client, err := genai.NewClient(ctx, &genai.ClientConfig{
		APIKey:  os.Getenv("GEMINI_API_KEY"),
		Backend: genai.BackendGeminiAPI,
	})
	if err != nil {
		log.Fatal(err)
	}

	parts := []*genai.Part{
		{Text: "What's this image about?"},
	}
	result, err := client.Models.GenerateContent(ctx, model, []*genai.Content{{Parts: parts}}, nil)
	if err != nil {
		log.Fatal(err)
	}

	return result
}
