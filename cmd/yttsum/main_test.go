package main

import (
	"flag"
	"os"
	"time"

	"github.com/bullshitsoftware/youtimetrack/internal/app"
	"github.com/bullshitsoftware/youtimetrack/internal/pkg/calendar"
	"github.com/bullshitsoftware/youtimetrack/internal/pkg/youtrack"
)

var (
	periodStart = time.Date(2007, 1, 1, 0, 0, 0, 0, time.UTC)
	periodEnd   = time.Date(2007, 1, 31, 23, 59, 59, 0, time.UTC)
	now1        = time.Date(2007, 1, 16, 3, 4, 5, 0, time.UTC)
	now2        = time.Date(2014, 1, 11, 3, 4, 5, 0, time.UTC)
)

type AppStub struct{}

func (s *AppStub) Load() {}

func (s *AppStub) NewCalendar() app.Calendar {
	return &CalendarStub{}
}

func (s *AppStub) NewYoutrack() app.Youtrack {
	return &YoutrackStub{}
}

func (yt *YoutrackStub) Delete(issueId, itemId string) error {
	panic("unexpected call")
}

type CalendarStub struct{}

func (s *CalendarStub) Period(now time.Time) *calendar.Period {
	if now.Equal(now1) || now.Equal(now2) {
		return &calendar.Period{
			Start: periodStart,
			End:   periodEnd,
		}
	}

	panic("enexpected now")
}

func (s *CalendarStub) Calc(start, end time.Time) int {
	if !start.Equal(periodStart) {
		panic("unexpected start " + start.String())
	}

	if end.Equal(now1) {
		return 100
	}

	if end.Equal(periodEnd) {
		return 230
	}

	panic("unexpected end")
}

type YoutrackStub struct{}

func (yt *YoutrackStub) WorkItems(start, end time.Time) ([]youtrack.WorkItem, error) {
	issue := &youtrack.Issue{IdReadable: "XY-123", Summary: "Do something cool"}
	items := []youtrack.WorkItem{
		{
			Issue:    issue,
			Date:     time.Date(2007, 1, 10, 3, 4, 5, 0, time.UTC).UnixMilli(),
			Duration: youtrack.Duration{Minutes: 30},
			Type:     youtrack.Type{Id: "123", Name: "Development"},
			Text:     "did something cool",
		},
		{
			Issue:    issue,
			Date:     time.Date(2007, 1, 15, 3, 4, 5, 0, time.UTC).UnixMilli(),
			Duration: youtrack.Duration{Minutes: 60},
			Type:     youtrack.Type{Id: "321", Name: "DevOps"},
			Text:     "opsed something cool",
		},
	}

	return items, nil
}

func (yt *YoutrackStub) WorkItemTypes() ([]youtrack.Type, error) {
	panic("unexpected call")
}

func (yt *YoutrackStub) Add(itemType youtrack.Type, issueId, duration, text string) error {
	panic("unexpected call")
}

func Example() {
	a = &AppStub{}
	now = now1
	main()

	flag.CommandLine = flag.NewFlagSet(os.Args[0], flag.ExitOnError)
	now = now2
	os.Args = []string{"yttsum", "-start", "2007-01-01", "-end", "2007-01-31"}
	main()

	// Output:
	// 1h30m / 1h40m / 3h50m (worked / today / month)
	// 1h30m / 3h50m (worked / month)
}
