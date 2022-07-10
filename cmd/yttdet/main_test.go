package main

import (
	"time"

	"github.com/bullshitsoftware/youtimetrack/internal/app"
	"github.com/bullshitsoftware/youtimetrack/internal/pkg/calendar"
	"github.com/bullshitsoftware/youtimetrack/internal/pkg/youtrack"
)

type AppStub struct{}

func (s *AppStub) Load() {}

func (s *AppStub) NewCalendar() (app.Calendar, error) {
	return &CalendarStub{}, nil
}

func (s *AppStub) NewYoutrack() app.Youtrack {
	return &YoutrackStub{}
}

type CalendarStub struct{}

func (s *CalendarStub) Period(now time.Time) *calendar.Period {
	if !now.Equal(time.Date(2007, 1, 2, 3, 4, 5, 0, time.UTC)) {
		panic("enexpected now")
	}

	return &calendar.Period{
		Start: time.Date(2007, 1, 1, 0, 0, 0, 0, time.UTC),
		End:   time.Date(2007, 1, 31, 23, 59, 59, 0, time.UTC),
	}
}

func (s *CalendarStub) Calc(start, end time.Time) int {
	panic("unexpected call")
}

type YoutrackStub struct{}

func (yt *YoutrackStub) WorkItems(start, end time.Time) ([]youtrack.WorkItem, error) {
	issue := &youtrack.Issue{IdReadable: "XY-123", Summary: "Do something cool"}
	items := []youtrack.WorkItem{
		{
			Id:       "110-12312",
			Issue:    issue,
			Date:     time.Date(2007, 1, 10, 3, 4, 5, 0, time.UTC).UnixMilli(),
			Duration: youtrack.Duration{Minutes: 30},
			Type:     youtrack.Type{Id: "123", Name: "Development"},
			Text:     "did something cool",
		},
		{
			Id:       "110-12313",
			Issue:    issue,
			Date:     time.Date(2007, 1, 15, 3, 4, 5, 0, time.UTC).UnixMilli(),
			Duration: youtrack.Duration{Minutes: 30},
			Type:     youtrack.Type{Id: "321", Name: "DevOps"},
			Text:     "opsed something\ncool",
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

func (yt *YoutrackStub) Delete(issueId, itemId string) error {
	panic("unexpected call")
}

func Example() {
	a = &AppStub{}
	now = time.Date(2007, 1, 2, 3, 4, 5, 0, time.UTC)
	main()

	// Output:
	// 2007-01-10	0h30m	XY-123	Do something cool
	// 110-12312		did something cool
	// 2007-01-15	0h30m	XY-123	Do something cool
	// 110-12313		opsed something
	// 			cool
}
