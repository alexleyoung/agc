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

func getService(ctx context.Context, userID string) (*calendar.Service, error) {
	client, err := auth.GetClient(userID)
	if err != nil {
		return nil, err
	}

	srv, err := calendar.NewService(ctx, option.WithHTTPClient(client))
	return srv, err
}
