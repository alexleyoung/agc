package calendar

import "google.golang.org/api/calendar/v3"

func CreateEvent(summary, description, start, end string) *calendar.Event {
	return &calendar.Event{
		Summary:     summary,
		Description: description,
		Start: &calendar.EventDateTime{
			DateTime: start,
		},
		End: &calendar.EventDateTime{
			DateTime: end,
		},
	}
}
