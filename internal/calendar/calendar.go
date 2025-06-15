package calendar

import (
	"log"
	"time"

	"github.com/alexleyoung/auto-gcal/internal/auth"
	"golang.org/x/net/context"
	"google.golang.org/api/calendar/v3"
	"google.golang.org/api/option"
)

func getService(ctx context.Context) (*calendar.Service, error) {
	client := auth.GetClient()

	srv, err := calendar.NewService(ctx, option.WithHTTPClient(client))
	return srv, err
}

func GetCalendarID(ctx context.Context, name string) (string, error) {
	srv, err := getService(ctx)
	if err != nil {
		log.Printf("Unable to retrieve calendar service: %v", err)
		return "", err
	}

	list, err := srv.CalendarList.List().Do()
	for _, cal := range list.Items {
		if cal.Summary == name {
			cal, err := srv.Calendars.Get(cal.Id).Do()
			if err != nil {
				log.Printf("Failed to retrieve calendar %s: %v", cal.Summary, err)
				return "", err
			}
			return cal.Id, nil
		}
	}

	return "primary", nil
}

// Get the current timestamp with the calendar's timezone
func Now() string {
	return time.Now().UTC().Format(time.RFC3339)
}

func CreateEvent(ctx context.Context, calendarID string, summary, description, start, end string) (*calendar.Event, error) {
	srv, err := getService(ctx)
	if err != nil {
		log.Printf("Unable to retrieve calendar service: %v", err)
		return &calendar.Event{}, err
	}

	if calendarID == "" {
		calendarID = "primary"
	}

	// get calendar's timezone
	cal, err := srv.Calendars.Get(calendarID).Do()
	if err != nil {
		log.Printf("Unable to retrieve calendar: %v", err)
		return &calendar.Event{}, err
	}

	ev := &calendar.Event{
		Summary:     summary,
		Description: description,
		Start: &calendar.EventDateTime{
			DateTime: start,
			TimeZone: cal.TimeZone,
		},
		End: &calendar.EventDateTime{
			DateTime: end,
			TimeZone: cal.TimeZone,
		},
	}

	ev, err = srv.Events.Insert(calendarID, ev).Do()
	if err != nil {
		log.Printf("Failed to create event \"%s\": %v", summary, err)
		return &calendar.Event{}, err
	}

	return ev, nil
}
