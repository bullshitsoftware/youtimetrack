package main

import (
	"fmt"
	"path"
	"time"
)

func Init(app *App) {
	app.SaveConfig()
	fmt.Println("Created " + path.Join(home, config))
}

func Summary(app *App) {
	start, end := app.Calendar.Period(now)
	month := app.Calendar.Calc(start, end)
	today := app.Calendar.Calc(start, now)
	items := app.Youtrack.Fetch(start, end)
	worked := 0
	for _, i := range items {
		worked += i.Duration.Minutes
	}
	fmt.Printf("%s / %s / %s (worked / today / month)\n",
		FormatMinutes(worked),
		FormatMinutes(today),
		FormatMinutes(month),
	)
}

func Details(app *App) {
	start, end := app.Calendar.Period(now)
	items := app.Youtrack.Fetch(start, end)
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
