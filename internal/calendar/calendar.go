package calendar

import (
	"log"

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

func CreateEvent(ctx context.Context, calendarID string, summary, description, start, end string) (*calendar.Event, error) {
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

	srv, err := getService(ctx)
	if err != nil {
		log.Printf("Unable to retrieve calendar service: %v", err)
		return &calendar.Event{}, err
	}

	ev, err = srv.Events.Insert(calendarID, ev).Do()
	if err != nil {
		log.Printf("Failed to create event \"%s\": %v", ev.Summary, err)
		return &calendar.Event{}, err
	}

	return ev, nil
}
