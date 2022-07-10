package main

import (
	"flag"
	"fmt"
	"os"
	"time"

	"github.com/bullshitsoftware/youtimetrack/internal/app"
)

func main() {
	a := app.Default()
	err := a.ReadConfig()
	app.ExitOnError(err)

	now := time.Now()
	period := a.Calendar.Period(now)
	fs := flag.NewFlagSet("summary", flag.ExitOnError)
	fs.Func("start", "start date (2006-01-02)", period.ParseStart)
	fs.Func("end", "end date (2006-01-02)", period.ParseEnd)
	fs.Parse(os.Args)

	month := a.Calendar.Calc(period.Start, period.End)
	items, err := a.Youtrack.WorkItems(period.Start, period.End)
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

	today := a.Calendar.Calc(period.Start, now)
	fmt.Printf("%s / %s / %s (worked / today / month)\n",
		app.FormatMinutes(worked),
		app.FormatMinutes(today),
		app.FormatMinutes(month),
	)
}
