package app

import (
	"os"
	"path"
	"time"

	cal "github.com/bullshitsoftware/youtimetrack/internal/pkg/calendar"
	yt "github.com/bullshitsoftware/youtimetrack/internal/pkg/youtrack"
)

const config = "config.json"

var home string

func init() {
	home = path.Join(os.Getenv("HOME"), ".config", "ytt")
}

type App struct {
	cfg Config
}

type Calendar interface {
	Period(now time.Time) *cal.Period
	Calc(start, end time.Time) int
}

type Youtrack interface {
	WorkItems(start, end time.Time) ([]yt.WorkItem, error)
	WorkItemTypes() ([]yt.Type, error)
	Add(itemType yt.Type, issueId, duration, text string) error
}

func (a *App) Load() {
	f, err := os.Open(path.Join(home, config))
	ExitOnError(err)
	defer f.Close()

	err = a.cfg.LoadJson(f)
	ExitOnError(err)
}

func (a *App) Save() string {
	err := os.MkdirAll(home, 0700)
	ExitOnError(err)

	p := path.Join(home, config)
	f, err := os.Create(p)
	ExitOnError(err)
	defer f.Close()

	err = a.cfg.SaveJson(f)
	ExitOnError(err)

	return p
}

func (a *App) NewCalendar() Calendar {
	c := a.cfg.Calendar
	return &cal.Calendar{
		DayDur:  c.DayDur,
		SDayDur: c.SDayDur,

		Weekends:  c.Weekends,
		Workdays:  c.Workdays,
		SWorkdays: c.SWorkdays,
		Holidays:  c.Holidays,
	}
}

func (a *App) NewYoutrack() Youtrack {
	c := a.cfg.Youtrack
	return &yt.Client{
		BaseUrl: c.BaseUrl,
		Token:   c.Token,
		Author:  c.Author,
	}
}
