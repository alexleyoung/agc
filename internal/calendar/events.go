package calendar

import (
	"context"
	"log"
	"time"

	"github.com/alexleyoung/agc/internal/types"
	"google.golang.org/api/calendar/v3"
)

func CreateEvent(ctx context.Context, session types.Session, calendarID, summary, description, start, end, timezone string) (*calendar.Event, error) {
	srv, err := getService(ctx)
	if err != nil {
		log.Printf("Unable to retrieve calendar service: %v", err)
		return &calendar.Event{}, err
	}

	if calendarID == "" {
		calendarID = "primary"
	}

	// get calendar's timezone
	if timezone == "" {
		cal, err := srv.Calendars.Get(calendarID).Do()
		if err != nil {
			log.Printf("Failed to fetch calendar %s: %v", calendarID, err)
			return &calendar.Event{}, err
		}
		timezone = cal.TimeZone
	}

	ev := &calendar.Event{
		Summary:     summary,
		Description: description,
		Start: &calendar.EventDateTime{
			DateTime: start,
			TimeZone: timezone,
		},
		End: &calendar.EventDateTime{
			DateTime: end,
			TimeZone: timezone,
		},
	}

	ev, err = srv.Events.Insert(calendarID, ev).Do()
	if err != nil {
		log.Printf("Failed to create event \"%s\": %v", summary, err)
		return &calendar.Event{}, err
	}

	return ev, nil
}

func GetEvents(ctx context.Context, session types.Session, calendarID string) ([]*calendar.Event, error) {
	list := make([]*calendar.Event, 0)

	srv, err := getService(ctx)
	if err != nil {
		log.Printf("Unable to retrieve calendar service: %v", err)
		return nil, err
	}

	time := time.Now().Format(time.RFC3339)
	orderBy := "startTime"
	maxResults := int64(50)
	events, err := srv.Events.List(calendarID).ShowDeleted(false).TimeMin(time).OrderBy(orderBy).MaxResults(maxResults).Do()
	if err != nil {
		log.Printf("Unable to retrieve events: %v", err)
		return nil, err
	}
	list = events.Items

	return list, nil
}
