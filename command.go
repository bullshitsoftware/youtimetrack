package main

import (
	"flag"
	"fmt"
	"path"
	"time"
)

func Init(app *App) {
	app.SaveConfig()
	fmt.Println("Created " + path.Join(home, config))
}

func Summary(app *App, args []string) {
	period := app.Calendar.Period(now)
	fs := flag.NewFlagSet("summary", flag.ExitOnError)
	fs.Func("start", "start date (2006-01-02)", period.ParseStart)
	fs.Func("end", "end date (2006-01-02)", period.ParseEnd)
	fs.Parse(args)

	month := app.Calendar.Calc(period.Start, period.End)
	items := app.Youtrack.Fetch(period.Start, period.End)
	worked := 0
	for _, i := range items {
		worked += i.Duration.Minutes
	}

	if now.Before(period.Start) || now.After(period.End) {
		fmt.Printf("%s / %s (worked / month)\n",
			FormatMinutes(worked),
			FormatMinutes(month),
		)

		return
	}

	today := app.Calendar.Calc(period.Start, now)
	fmt.Printf("%s / %s / %s (worked / today / month)\n",
		FormatMinutes(worked),
		FormatMinutes(today),
		FormatMinutes(month),
	)
}

func Details(app *App, args []string) {
	period := app.Calendar.Period(now)
	fs := flag.NewFlagSet("details", flag.ExitOnError)
	fs.Func("start", "start date (2006-01-02)", period.ParseStart)
	fs.Func("end", "end date (2006-01-02)", period.ParseEnd)
	fs.Parse(args)

	items := app.Youtrack.Fetch(period.Start, period.End)
	for _, i := range items {
		date := time.Unix(i.Date/1000, 0)
		fmt.Println(
			date.Format("2006-01-02"),
			FormatMinutes(i.Duration.Minutes),
			i.Issue.IdReadable,
			i.Issue.Summary,
		)
	}
}
