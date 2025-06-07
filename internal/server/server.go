package server

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/alexleyoung/auto-gcal/internal/ai"
	"github.com/alexleyoung/auto-gcal/internal/auth"
	"google.golang.org/api/calendar/v3"
	"google.golang.org/api/option"
)

func Run() {
	InitConfig()

	ctx := context.Background()
	client := auth.GetClient()

	srv, err := calendar.NewService(ctx, option.WithHTTPClient(client))
	if err != nil {
		log.Fatalf("Unable to retrieve Calendar client: %v", err)
	}

	calendarList, err := srv.CalendarList.List().Do()
	for _, cal := range calendarList.Items {
		fmt.Println(cal.Id)
	}

	t := time.Now().Format(time.RFC3339)

	events, err := srv.Events.List("primary").ShowDeleted(false).
		SingleEvents(true).TimeMin(t).MaxResults(10).OrderBy("startTime").Do()
	if err != nil {
		log.Fatalf("Unable to retrieve next ten of the user's events: %v", err)
	}
	fmt.Println("Upcoming events:")
	if len(events.Items) == 0 {
		fmt.Println("No upcoming events found.")
	} else {
		for _, item := range events.Items {
			date := item.Start.DateTime
			if date == "" {
				date = item.Start.Date
			}
			fmt.Printf("%v (%v)\n", item.Summary, date)
		}
	}

	ai.Chat(ctx, "Hey Gemini")
}
