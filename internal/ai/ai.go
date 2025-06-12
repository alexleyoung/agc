package ai

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"

	"github.com/alexleyoung/auto-gcal/internal/calendar"
	"google.golang.org/genai"
)

const (
	MAX_STEPS     = 5
	SYSTEM_PROMPT = `You are an intelligent assistant that helps users manage their Google Calendar.
Your job is to extract relevant details from user input—like the event title, start and end times, and optional descriptions—and call the appropriate function to schedule the event.
When the user describes an event, respond only by calling the create_event function with the appropriate parameters.
After scheduling, confirm success by summarizing the event details back to the user in natural language.
If any required information is missing or ambiguous, ask the user for clarification. Do not guess.`
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
			"start":       {Type: "string", Description: "The time, as a combined date-time value (formatted according to RFC3339). Required."},
			"end":         {Type: "string", Description: "The time, as a combined date-time value (formatted according to RFC3339). Required"},
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
	for step := 0; step < MAX_STEPS; step++ {
		// prompt model
		result, err = client.Models.GenerateContent(ctx, model, history, &genai.GenerateContentConfig{
			SystemInstruction: genai.NewContentFromText(SYSTEM_PROMPT, genai.RoleUser),
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
			argsJSON, _ := json.Marshal(fn.Args)
			log.Printf("Model requested function: %s\nWith args: %s", fn.Name, argsJSON)

			out, err := executeFunctionCall(ctx, fn)
			if err != nil {
				log.Printf("Error executing function: %v", err)
				return &genai.GenerateContentResponse{}, err
			}
			log.Print("function execution successful")

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

		data, err := json.Marshal(fn.Args)
		if err != nil {
			return "", fmt.Errorf("failed to marshal args: %w", err)
		}
		if err := json.Unmarshal(data, &args); err != nil {
			return "", fmt.Errorf("failed to decode args into struct: %w", err)
		}

		ev, err := calendar.CreateEvent(ctx, args.CalendarID, args.Summary, args.Description, args.Start, args.End)
		if err != nil {
			return "", err
		}
		return fmt.Sprintf("Successfully created event \"%s\"", ev.Summary), nil
	}
	return "", fmt.Errorf("Unknown function: %s", fn.Name)
}
