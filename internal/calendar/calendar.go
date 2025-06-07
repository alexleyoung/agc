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

func GetCalendar(ctx context.Context, name string) (*calendar.Calendar, error) {
	srv, err := getService(ctx)
	if err != nil {
		log.Printf("Unable to retrieve calendar service: %v", err)
		return &calendar.Calendar{}, err
	}

	list, err := srv.CalendarList.List().Do()
	for _, cal := range list.Items {
		if cal.Summary == name {
			return srv.Calendars.Get(cal.Id).Do()
		}
	}

	return srv.Calendars.Get("primary").Do()
}

func CreateEvent(ctx context.Context, cal *calendar.CalendarListEntry, summary, description, start, end string) (*calendar.Event, error) {
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

	ev, err = srv.Events.Insert(cal.Id, ev).Do()
	if err != nil {
		log.Printf("Failed to create event \"%s\": %v", ev.Summary, err)
		return &calendar.Event{}, err
	}

	return ev, nil
}
