package main

import (
	"context"
	"fmt"

	"github.com/alexleyoung/auto-gcal/internal/calendar"
	"github.com/alexleyoung/auto-gcal/internal/server"
)

func main() {
	time, err := calendar.GetUserCurrentDateTime(context.Background(), "primary")
	if err != nil {
		return
	}
	fmt.Print(time)
	server.Run()
}
