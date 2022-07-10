package main

import (
	"flag"
	"fmt"
	"time"

	"github.com/bullshitsoftware/youtimetrack/internal/app"
)

type App interface {
	Load()
	NewCalendar() (app.Calendar, error)
	NewYoutrack() app.Youtrack
}

var (
	a   App       = &app.App{}
	now time.Time = time.Now()
)

func main() {
	a.Load()

	cal, err := a.NewCalendar()
	app.ExitOnError(err)
	yt := a.NewYoutrack()

	period := cal.Period(now)
	flag.Func("start", "start date (2006-01-02)", period.ParseStart)
	flag.Func("end", "end date (2006-01-02)", period.ParseEnd)
	flag.Parse()

	month := cal.Calc(period.Start, period.End)
	items, err := yt.WorkItems(period.Start, period.End)
	app.ExitOnError(err)
	worked := 0
	for _, i := range items {
		worked += i.Duration.Minutes
	}

	if now.Before(period.Start) || now.After(period.End) {
		fmt.Printf("%s / %s (worked / month)\n",
			app.FormatMinutes(worked),
			app.FormatMinutes(month),
		)

		return
	}

	today := cal.Calc(period.Start, now)
	fmt.Printf("%s / %s / %s (worked / today / month)\n",
		app.FormatMinutes(worked),
		app.FormatMinutes(today),
		app.FormatMinutes(month),
	)
}
