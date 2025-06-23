package calendar

import (
	"log"
	"time"

	"github.com/alexleyoung/auto-gcal/internal/auth"
	"golang.org/x/net/context"
	"google.golang.org/api/calendar/v3"
	"google.golang.org/api/option"
)

// Get the current UTC timestamp as string
func Now() string {
	return time.Now().UTC().Format(time.RFC3339)
}

func GetCalendarID(ctx context.Context, userID, name string) (string, error) {
	srv, err := getService(ctx, userID)
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

func CreateEvent(ctx context.Context, userID, calendarID, summary, description, start, end, timezone string) (*calendar.Event, error) {
	srv, err := getService(ctx, userID)
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

func GetEvents(ctx context.Context, userID, calendarID string) ([]*calendar.Event, error) {
	list := make([]*calendar.Event, 0)

	srv, err := getService(ctx, userID)
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

func getService(ctx context.Context, userID string) (*calendar.Service, error) {
	client, err := auth.GetClient(userID)
	if err != nil {
		return nil, err
	}

	srv, err := calendar.NewService(ctx, option.WithHTTPClient(client))
	return srv, err
}
