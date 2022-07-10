package main

import (
	"flag"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/bullshitsoftware/youtimetrack/internal/app"
)

func main() {
	a := app.Default()
	err := a.ReadConfig()
	app.ExitOnError(err)

	now := time.Now()
	period := a.Calendar.Period(now)
	fs := flag.NewFlagSet("details", flag.ExitOnError)
	fs.Func("start", "start date (2006-01-02)", period.ParseStart)
	fs.Func("end", "end date (2006-01-02)", period.ParseEnd)
	fs.Parse(os.Args)

	items, err := a.Youtrack.WorkItems(period.Start, period.End)
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
