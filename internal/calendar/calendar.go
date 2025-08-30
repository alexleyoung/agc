package calendar

import (
	"context"
	"log"

	"github.com/alexleyoung/agc/internal/auth"
	"github.com/alexleyoung/agc/internal/types"
	"google.golang.org/api/calendar/v3"
	"google.golang.org/api/option"
)

func GetCalendar(ctx context.Context, session types.Session, calendar string) (*calendar.Calendar, error) {
	srv, err := getService(ctx)
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

func getService(ctx context.Context) (*calendar.Service, error) {
	client := auth.GetClient()
	srv, err := calendar.NewService(ctx, option.WithHTTPClient(client))
	return srv, err
}
