package llm

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
	{
		Name:        "get_events",
		Description: "Fetches events from a calendar",
		Parameters: &genai.Schema{
			Type: "object",
			Properties: map[string]*genai.Schema{
				"calendar_id": {Type: "string", Description: "The ID of the calendar to fetch events from. Required."},
				"min_time":    {Type: "string", Description: "The minimum time to fetch events from. Optional."},
				"max_time":    {Type: "string", Description: "The maximum time to fetch events from. Optional."},
				"max_results": {Type: "number", Description: "The maximum number of events to fetch. Optional."},
			},
		},
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

func get_events_call(ctx context.Context, argsJSON []byte) (string, error) {
	var args struct {
		CalendarID string `json:"calendar_id"`
		MinTime    string `json:"min_time"`
		MaxTime    string `json:"max_time"`
		MaxResults int64  `json:"max_results"`
	}

	if err := json.Unmarshal(argsJSON, &args); err != nil {
		return "", fmt.Errorf("failed to decode args into struct: %w", err)
	}

	events, err := calendar.GetEvents(ctx, args.CalendarID, args.MinTime, args.MaxTime, args.MaxResults)
	if err != nil {
		return "", err
	}

	var event_strings string
	for _, event := range events {
		event_strings += formatEvent(event)
		event_strings += "\n"
	}
	return fmt.Sprintf("Events:\n%s", event_strings), nil
}

func executeFunctionCall(ctx context.Context, name string, argsJSON []byte) (string, error) {
	switch name {
	case "create_event":
		return create_event_call(ctx, argsJSON)
	case "quick_add_event":
		return quick_add_event_call(ctx, argsJSON)
	case "list_calendars":
		return list_calendars_call(ctx)
	case "get_current_time":
		return get_current_time_call(ctx)
	case "get_events":
		return get_events_call(ctx, argsJSON)
	}
	return "", fmt.Errorf("Unknown function: %s", name)
}
