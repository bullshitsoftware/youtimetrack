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
	start, end := app.Calendar.Period(now)
	fs := flag.NewFlagSet("summary", flag.ExitOnError)
	fs.Func("start", "start date (2006-01-02)", func(s string) error {
		var err error
		start, err = time.Parse("2006-01-02", s)
		if err != nil {
			fmt.Printf("%v\n", start)
			return err
		}
		return nil
	})
	fs.Func("end", "end date (2006-01-02)", func(s string) error {
		var err error
		end, err = time.Parse("2006-01-02", s)
		if err != nil {
			return err
		}
		end = time.Date(end.Year(), end.Month(), end.Day(), 23, 59, 59, 0, time.UTC)
		return nil
	})
	fs.Parse(args)

	month := app.Calendar.Calc(start, end)
	items := app.Youtrack.Fetch(start, end)
	worked := 0
	for _, i := range items {
		worked += i.Duration.Minutes
	}

	if now.Before(start) || now.After(end) {
		fmt.Printf("%s / %s (worked / month)\n",
			FormatMinutes(worked),
			FormatMinutes(month),
		)

		return
	}

	today := app.Calendar.Calc(start, now)
	fmt.Printf("%s / %s / %s (worked / today / month)\n",
		FormatMinutes(worked),
		FormatMinutes(today),
		FormatMinutes(month),
	)
}

func Details(app *App, args []string) {
	start, end := app.Calendar.Period(now)
	fs := flag.NewFlagSet("details", flag.ExitOnError)
	fs.Func("start", "start date (2006-01-02)", func(s string) error {
		var err error
		start, err = time.Parse("2006-01-02", s)
		if err != nil {
			fmt.Printf("%v\n", start)
			return err
		}
		return nil
	})
	fs.Func("end", "end date (2006-01-02)", func(s string) error {
		var err error
		end, err = time.Parse("2006-01-02", s)
		if err != nil {
			return err
		}
		end = time.Date(end.Year(), end.Month(), end.Day(), 23, 59, 59, 0, time.UTC)
		return nil
	})
	fs.Parse(args)

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
