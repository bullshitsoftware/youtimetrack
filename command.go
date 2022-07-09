package main

import (
	"flag"
	"fmt"
	"path"
	"strings"
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
		fmt.Printf(
			"%s\t%s\t%s\t%s\n",
			date.Format("2006-01-02"),
			FormatMinutes(i.Duration.Minutes),
			i.Issue.IdReadable,
			i.Issue.Summary,
		)
		for _, s := range strings.Split(i.Text, "\n") {
			fmt.Printf("\t\t\t%s\n", s)
		}
	}
}

func Add(app *App, args []string) {
	if len(args) != 4 {
		panic("Invalid arguments number")
	}
	typeName := strings.ToLower(args[0])
	types := app.Youtrack.WorkItemTypes()
	var t Type
	for _, i := range types {
		s := strings.ToLower(i.Name)
		if strings.HasPrefix(s, typeName) {
			t = i
			break
		}
	}
	issue := args[1]
	duration := args[2]
	text := args[3]

	app.Youtrack.Add(t, issue, duration, text)
}
