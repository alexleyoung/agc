package calendar

import (
	"context"
	"log"
	"time"

	"google.golang.org/api/calendar/v3"
)

func QuickAddEvent(ctx context.Context, calendarID, query string) (*calendar.Event, error) {
	srv, err := getService(ctx)
	if err != nil {
		log.Printf("Unable to retrieve calendar service: %v", err)
		return nil, err
	}

	ev, err := srv.Events.QuickAdd(calendarID, query).Do()
	if err != nil {
		log.Printf("Unable to retrieve event: %v", err)
		return nil, err
	}

	return ev, nil
}

func CreateEvent(ctx context.Context, calendarID, summary, description, start, end, timezone string) (*calendar.Event, error) {
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
		Recurrence: nil,
	}

	ev, err = srv.Events.Insert(calendarID, ev).Do()
	if err != nil {
		log.Printf("Failed to create event \"%s\": %v", summary, err)
		return &calendar.Event{}, err
	}

	return ev, nil
}

func UpdateEvent(ctx context.Context, calendarID string, event *calendar.Event) error {
	srv, err := getService(ctx)
	if err != nil {
		log.Printf("Unable to retrieve calendar service: %v", err)
		return err
	}
	_, err = srv.Events.Update(calendarID, event.Id, event).Do()
	return err
}

func GetEvent(ctx context.Context, calendarID, eventID string) (*calendar.Event, error) {
	srv, err := getService(ctx)
	if err != nil {
		log.Printf("Unable to retrieve calendar service: %v", err)
		return nil, err
	}

	event, err := srv.Events.Get(calendarID, eventID).Do()
	if err != nil {
		log.Printf("Unable to retrieve event: %v", err)
		return nil, err
	}

	return event, nil
}

func GetEvents(ctx context.Context, calendarID, minTime, maxTime string, maxResults int64) ([]*calendar.Event, error) {
	srv, err := getService(ctx)
	if err != nil {
		log.Printf("Unable to retrieve calendar service: %v", err)
		return nil, err
	}

	if minTime == "" {
		minTime = time.Now().Format(time.RFC3339)
	}
	query := srv.Events.List(calendarID).ShowDeleted(false).OrderBy("startTime").MaxResults(maxResults).TimeMin(minTime)
	if maxTime != "" {
		query = query.TimeMax(maxTime)
	}

	events, err := query.Do()
	if err != nil {
		log.Printf("Unable to retrieve events: %v", err)
		return nil, err
	}

	return events.Items, nil
}
