package ai

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/alexleyoung/agc/internal/calendar"
	"google.golang.org/genai"
)

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

func create_event_call(ctx context.Context, argsJSON []byte) (string, error) {
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
	return fmt.Sprintf("Successfully created event \"%s\"", formatEvent(ev)), nil
}

func quick_add_event_call(ctx context.Context, argsJSON []byte) (string, error) {
	var args struct {
		CalendarID string `json:"calendar_id"`
		Query      string `json:"query"`
	}

	if err := json.Unmarshal(argsJSON, &args); err != nil {
		return "", fmt.Errorf("failed to decode args into struct: %w", err)
	}

	ev, err := calendar.QuickAddEvent(ctx, args.CalendarID, args.Query)
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("Successfully created event \"%s\"", formatEvent(ev)), nil
}

func list_calendars_call(ctx context.Context) (string, error) {
	cals, err := calendar.ListCalendars(ctx)
	if err != nil {
		return "", err
	}
	var cal_strings string
	for _, cal := range cals {
		cal_strings += formatCalendar(cal)
		cal_strings += "\n"
	}
	return fmt.Sprintf("Calendars:\n%s", cal_strings), nil
}

func get_current_time_call(ctx context.Context) (string, error) {
	time := time.Now().UTC().Format(time.RFC3339)
	return fmt.Sprintf("Current time: %s", time), nil
}
