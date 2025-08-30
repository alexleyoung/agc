package ai

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/alexleyoung/agc/internal/calendar"
	"github.com/spf13/viper"

	"google.golang.org/genai"
)

const (
	MAX_STEPS     = 10
	SYSTEM_PROMPT = `You are a helpful assistant that enables users to interact with Google Calendar through natural language.
Your role is to interpret user requests and translate them into appropriate calendar actions.

Capabilities
Event management: create, update, delete, and retrieve events.
Information queries: check schedules, availability, upcoming events, event details, and reminders.
Scheduling assistance: suggest times, handle recurring events, and resolve conflicts.
Resource awareness: understand Google Calendar entities such as events, attendees, reminders, time zones, recurrence rules, and resources.

Instructions
Always clarify ambiguities before making changes (e.g., confirm times, dates, and event names if unclear).
Respect user intent precisely—never add, modify, or delete events without explicit instruction.
When answering, provide both a natural language response and a structured action representation (e.g., API call, JSON payload, or step summary, depending on the integration).
Keep answers concise and user-friendly while surfacing important details (time, date, participants).
Always account for time zones, recurring rules, and shared calendars when relevant.
If a request cannot be fulfilled with Google Calendar (e.g., booking a restaurant), politely explain and suggest alternatives.

Style
Be conversational, concise, and professional.
Use confirmation questions when multiple interpretations are possible.
When retrieving data, provide summaries instead of raw dumps unless the user requests full details.`
)

var TEMPERATURE float32 = 0

var functionDeclarations = []*genai.FunctionDeclaration{
	{
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
				"timezone":    {Type: "string", Description: "The timezone the datetime represents. Required."},
			},
		},
	},
	{
		Name:        "quick_add_event",
		Description: "Creates a new event in the user's calendar with natural language.",
		Parameters: &genai.Schema{
			Type: "object",
			Properties: map[string]*genai.Schema{
				"calendar_id": {Type: "string", Description: "The ID of the calendar to create the event in. Required."},
				"query":       {Type: "string", Description: "The query to use to create the event. Required."},
			},
		},
	},
	{
		Name:        "list_calendars",
		Description: "Lists all of the user's calendars",
	},
	{
		Name:        "get_current_time",
		Description: "Fetches the current time in UTC as an RFC3339 string",
	},
}

func Chat(ctx context.Context, model string, history []*genai.Content, prompt string) (*genai.GenerateContentResponse, error) {
	client, err := genai.NewClient(ctx, &genai.ClientConfig{
		APIKey:  viper.Get("GEMINI_API_KEY").(string),
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

			out, err := executeFunctionCall(ctx, fn.Name, args)
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

func executeFunctionCall(ctx context.Context, name string, argsJSON []byte) (string, error) {
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

		ev, err := calendar.CreateEvent(ctx, args.CalendarID, args.Summary, args.Description, args.Start, args.End, args.Timezone)
		if err != nil {
			return "", err
		}
		return fmt.Sprintf("Successfully created event \"%s\"", ev.Summary), nil

	case "list_calendars":
		cals, err := calendar.ListCalendars(ctx)
		if err != nil {
			return "", err
		}
		var cal_strings string
		for _, cal := range cals {
			cal_strings += fmt.Sprintf("Summary: %s\n", cal.Summary)
			cal_strings += fmt.Sprintf("ID: %s\n", cal.Id)
			cal_strings += fmt.Sprintf("Description: %s\n", cal.Description)
			cal_strings += fmt.Sprintf("Timezone: %s\n", cal.TimeZone)
			cal_strings += fmt.Sprintf("IsPrimary: %s\n", cal.Primary)
			cal_strings += "\n"
		}
		return fmt.Sprintf("Calendars:\n%s", cal_strings), nil

	case "get_current_time":
		time := time.Now().UTC().Format(time.RFC3339)
		return fmt.Sprintf("Current time: %s", time), nil
	}

	return "", fmt.Errorf("Unknown function: %s", name)
}
