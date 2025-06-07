package calendar

import (
	"log"

	"github.com/alexleyoung/auto-gcal/internal/auth"
	"golang.org/x/net/context"
	"google.golang.org/api/calendar/v3"
	"google.golang.org/api/option"
)

func CreateEvent(ctx context.Context, calendarID, summary, description, start, end string) (*calendar.Event, error) {
	ev := &calendar.Event{
		Summary:     summary,
		Description: description,
		Start: &calendar.EventDateTime{
			DateTime: start,
		},
		End: &calendar.EventDateTime{
			DateTime: end,
		},
	}

	client := auth.GetClient()
	srv, err := calendar.NewService(ctx, option.WithHTTPClient(client))
	if err != nil {
		log.Fatalf("Unable to retrieve Calendar client: %v", err)
		return &calendar.Event{}, err
	}

	ev, err = srv.Events.Insert(calendarID, ev).Do()
	if err != nil {
		log.Fatalf("Failed to create event \"%s\": %v", ev.Summary, err)
		return &calendar.Event{}, err
	}

	return ev, nil
}
