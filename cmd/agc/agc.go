package main

import (
	"context"
	"fmt"

	"github.com/alexleyoung/agc/internal/calendar"
	"github.com/alexleyoung/agc/internal/cli"
)

func main() {
	cals, _ := calendar.ListCalendars(context.Background())
	var cal_strings string
	for _, cal := range cals {
		cal_strings += fmt.Sprintf("%s\n", cal.Summary)
	}
	cli.Execute()
}
