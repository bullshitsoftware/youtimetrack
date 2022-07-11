package app

import (
	"encoding/json"
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
	Delete(issueId, itemId string) error
}

func (a *App) Load() {
	f, err := os.Open(path.Join(home, config))
	ExitOnError(err)
	defer f.Close()

	err = json.NewDecoder(f).Decode(&a.cfg)
	ExitOnError(err)
}

func (a *App) Save() string {
	err := os.MkdirAll(home, 0700)
	ExitOnError(err)

	p := path.Join(home, config)
	f, err := os.Create(p)
	ExitOnError(err)
	defer f.Close()

	enc := json.NewEncoder(f)
	enc.SetIndent("", "  ")
	err = enc.Encode(&a.cfg)
	ExitOnError(err)

	return p
}

func (a *App) NewCalendar() (Calendar, error) {
	c := a.cfg.Calendar
	holidays := []string{}
	for _, v := range c.Vacations {
		s, err := time.Parse("2006-01-02", v.Start)
		if err != nil {
			return nil, err
		}
		e, err := time.Parse("2006-01-02", v.End)
		if err != nil {
			return nil, err
		}
		for cur := s; cur.Before(e) || cur.Equal(e); cur = cur.AddDate(0, 0, 1) {
			holidays = append(holidays, cur.Format("2006-01-02"))
		}
	}

	holidays = append(holidays, c.Holidays...)
	calendar := &cal.Calendar{
		DayDur:  c.DayDur,
		SDayDur: c.SDayDur,

		Weekends:  c.Weekends,
		Workdays:  c.Workdays,
		SWorkdays: c.SWorkdays,
		Holidays:  holidays,
	}

	return calendar, nil
}

func (a *App) NewYoutrack() Youtrack {
	c := a.cfg.Youtrack
	return &yt.Client{
		BaseUrl: c.BaseUrl,
		Token:   c.Token,
		Author:  c.Author,
	}
}
