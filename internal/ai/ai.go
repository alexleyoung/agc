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
	MAX_STEPS     = 10
	SYSTEM_PROMPT = `You are an intelligent assistant that helps users manage their Google Calendar.
Your job is to extract relevant details from user input—like the event title, start and end times, and optional descriptions—and call the appropriate functions to schedule the event.
Always try to schedule the event; if it fails, simply let the user know the error.
After scheduling, confirm success by summarizing the event details back to the user in natural language.`
)

var TEMPERATURE float32 = 0

var functionDeclarations = []*genai.FunctionDeclaration{{
	Name:        "create_event",
	Description: "Creates a new event in the user's calendar.",
	Parameters: &genai.Schema{
		Type: "object",
		Properties: map[string]*genai.Schema{
			"calendar_id": {Type: "string", Description: "The ID of the calendar to create the event in. Optional."},
			"summary":     {Type: "string", Description: "The title of the event. Required."},
			"description": {Type: "string", Description: "The description of the event. Optional."},
			"start":       {Type: "string", Description: "The time, as a combined date-time value (formatted according to RFC3339) with NO offset. Required."},
			"end":         {Type: "string", Description: "The time, as a combined date-time value (formatted according to RFC3339) with NO offset. Required."},
			"timezone":    {Type: "string", Description: "The timezone the datetime represents. Optional."},
		},
	},
},
	{
		Name:        "get_current_time",
		Description: "Fetches the current time in UTC as an RFC3339 string",
	},
}

func Chat(ctx context.Context, session types.Session, model string, history []*genai.Content, prompt string) (*genai.GenerateContentResponse, error) {
	client, err := genai.NewClient(ctx, &genai.ClientConfig{
		APIKey:  os.Getenv("GEMINI_API_KEY"),
		Backend: genai.BackendGeminiAPI,
	})
	if err != nil {
		log.Printf("Error generating AI client: %v", err)
		return &genai.GenerateContentResponse{}, err
	}

	history = append(history, &genai.Content{
		Role:  "user",
		Parts: []*genai.Part{{Text: prompt}},
	})

	var result *genai.GenerateContentResponse
	for range MAX_STEPS {
		// prompt model
		result, err = client.Models.GenerateContent(ctx, model, history, &genai.GenerateContentConfig{
			SystemInstruction: genai.NewContentFromText(SYSTEM_PROMPT, genai.RoleUser),
			Tools: []*genai.Tool{
				{FunctionDeclarations: functionDeclarations},
			},
			Temperature: &TEMPERATURE,
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

			args, err := json.Marshal(fn.Args)
			if err != nil {
				return &genai.GenerateContentResponse{}, fmt.Errorf("failed to marshal args: %w", err)
			}
			log.Printf("Model requested function: %s\nWith args: %s", fn.Name, args)

			out, err := executeFunctionCall(ctx, session, fn.Name, args)
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

func executeFunctionCall(ctx context.Context, session types.Session, name string, argsJSON []byte) (string, error) {
	switch name {
	case "create_event":
		var args struct {
			CalendarID  string `json:"calendar_id"`
			Summary     string `json:"summary"`
			Description string `json:"description"`
			Start       string `json:"start"`
			End         string `json:"end"`
			Timezone    string `json:"timezone"`
		}

		if err := json.Unmarshal(argsJSON, &args); err != nil {
			return "", fmt.Errorf("failed to decode args into struct: %w", err)
		}

		ev, err := calendar.CreateEvent(ctx, session, args.CalendarID, args.Summary, args.Description, args.Start, args.End, args.Timezone)
		if err != nil {
			return "", err
		}
		return fmt.Sprintf("Successfully created event \"%s\"", ev.Summary), nil

	case "get_current_time":
		time := calendar.Now()
		return fmt.Sprintf("Current time: %s", time), nil
	}

	return "", fmt.Errorf("Unknown function: %s", name)
}
