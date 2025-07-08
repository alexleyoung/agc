package calendar

import (
	"log"

	"github.com/alexleyoung/auto-gcal/internal/auth"
	"github.com/alexleyoung/auto-gcal/internal/types"
	"golang.org/x/net/context"
	"google.golang.org/api/calendar/v3"
	"google.golang.org/api/option"
)

func GetCalendar(ctx context.Context, session types.Session, calendar string) (*calendar.Calendar, error) {
	srv, err := getService(ctx, session)
	if err != nil {
		log.Printf("Unable to retrieve calendar service: %v", err)
		return nil, err
	}

	cal, err := srv.Calendars.Get(calendar).Do()
	if err != nil {
		log.Printf("Unable to retrieve calendar: %v", err)
		return nil, err
	}

	return cal, nil
}

func getService(ctx context.Context, session types.Session) (*calendar.Service, error) {
	client, err := auth.GetClient(session)
	if err != nil {
		return nil, err
	}

	srv, err := calendar.NewService(ctx, option.WithHTTPClient(client))
	return srv, err
}
