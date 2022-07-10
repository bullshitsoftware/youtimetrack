package main

import (
	"flag"
	"fmt"
	"strings"
	"time"

	"github.com/bullshitsoftware/youtimetrack/internal/app"
)

type App interface {
	Load()
	NewCalendar() app.Calendar
	NewYoutrack() app.Youtrack
}

var (
	a   App       = &app.App{}
	now time.Time = time.Now()
)

func main() {
	a.Load()

	cal := a.NewCalendar()
	yt := a.NewYoutrack()

	period := cal.Period(now)
	flag.Func("start", "start date (2006-01-02)", period.ParseStart)
	flag.Func("end", "end date (2006-01-02)", period.ParseEnd)
	flag.Parse()

	items, err := yt.WorkItems(period.Start, period.End)
	app.ExitOnError(err)
	for _, i := range items {
		date := time.Unix(i.Date/1000, 0)
		fmt.Printf(
			"%s\t%s\t%s\t%s\n",
			date.Format("2006-01-02"),
			app.FormatMinutes(i.Duration.Minutes),
			i.Issue.IdReadable,
			i.Issue.Summary,
		)
		for _, s := range strings.Split(i.Text, "\n") {
			fmt.Printf("\t\t\t%s\n", s)
		}
	}
}
