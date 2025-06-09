package ai

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/alexleyoung/auto-gcal/internal/calendar"
	"google.golang.org/genai"
)

const (
	maxSteps = 5
)

var functionDeclarations = []*genai.FunctionDeclaration{{
	Name:        "create_event",
	Description: "Creates a new event in the user's calendar.",
	Parameters: &genai.Schema{
		Type: "object",
		Properties: map[string]*genai.Schema{
			"calendar_id": {Type: "string", Description: "The ID of the calendar to create the event in. Default to \"primary\"."},
			"summary":     {Type: "string", Description: "The title of the event. Required."},
			"description": {Type: "string", Description: "The description of the event. Default to empty."},
			"start":       {Type: "string", Description: "The time, as a combined date-time value (formatted according to RFC3339) with a timezone offset. Required."},
			"end":         {Type: "string", Description: "The time, as a combined date-time value (formatted according to RFC3339) with a timezone offset. Required"},
		},
	},
}}

func Chat(ctx context.Context, model string, prompt string) (*genai.GenerateContentResponse, error) {
	client, err := genai.NewClient(ctx, &genai.ClientConfig{
		APIKey:  os.Getenv("GEMINI_API_KEY"),
		Backend: genai.BackendGeminiAPI,
	})
	if err != nil {
		log.Printf("Error generating AI client: %v", err)
		return &genai.GenerateContentResponse{}, err
	}

	history := []*genai.Content{
		{Role: "user", Parts: []*genai.Part{{Text: prompt}}},
	}
	var result *genai.GenerateContentResponse
	for step := 0; step < maxSteps; step++ {
		// prompt model
		result, err = client.Models.GenerateContent(ctx, model, history, &genai.GenerateContentConfig{
			Tools: []*genai.Tool{
				&genai.Tool{
					FunctionDeclarations: functionDeclarations,
				},
			},
		})
		if err != nil {
			log.Printf("Error getting model response: %v", err)
			return &genai.GenerateContentResponse{}, err
		}

		history = append(history, result.Candidates[0].Content)

		// check for function calls
		fns := result.FunctionCalls()
		if len(fns) > 0 {
			fn := fns[0]
			log.Printf("Model requested function: %s\nWith args: %s", fn.Name, fn.Args)

			out, err := executeFunctionCall(ctx, fn)
			if err != nil {
				log.Printf("Error executing function: %v", err)
				return &genai.GenerateContentResponse{}, err
			}

			history = append(history, &genai.Content{
				Role:  "model",
				Parts: []*genai.Part{{Text: fn.Name + " results: "}, {Text: out}},
			})
			continue
		}
		return result, nil
	}
	return &genai.GenerateContentResponse{}, fmt.Errorf("Max steps reached without resolution")
}

func executeFunctionCall(ctx context.Context, fn *genai.FunctionCall) (string, error) {
	switch fn.Name {
	case "create_event":
		var args struct {
			CalendarID  string `json:"calendar_id"`
			Summary     string `json:"summary"`
			Description string `json:"description"`
			Start       string `json:"start"`
			End         string `json:"end"`
		}
		_, err := calendar.CreateEvent(ctx, args.CalendarID, args.Summary, args.Description, args.Start, args.End)
		if err != nil {
			return "", err
		}
		break
	}
	return "", fmt.Errorf("Unknown function: %s", fn.Name)
}
