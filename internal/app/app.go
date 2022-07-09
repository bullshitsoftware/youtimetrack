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

type App struct {
	Youtrack yt.Client    `json:"youtrack"`
	Calendar cal.Calendar `json:"calendar"`
}

func Default() *App {
	return &App{
		yt.Client{
			BaseUrl: "http://localhost:2378/api",
			Token:   "your-token",
			Author:  "your-user-uuid",
		},
		cal.Calendar{
			DayDur:    8 * 60,
			SDayDur:   7 * 60,
			Weekends:  []time.Weekday{time.Saturday, time.Sunday},
			Workdays:  []string{"2022-02-05"},
			SWorkdays: []string{"2022-02-21"},
			Holidays:  []string{"2022-02-22"},
		},
	}
}

func (a *App) SaveConfig(home string) string {
	err := os.MkdirAll(home, 0700)
	if err != nil {
		panic(err)
	}

	p := path.Join(home, config)
	f, err := os.Create(p)
	if err != nil {
		panic(err)
	}
	defer f.Close()

	enc := json.NewEncoder(f)
	enc.SetIndent("", "  ")
	err = enc.Encode(a)
	if err != nil {
		panic(err)
	}

	return p
}

func (a *App) ReadConfig(home string) {
	f, err := os.Open(path.Join(home, config))
	if err != nil {
		panic(err)
	}
	defer f.Close()

	err = json.NewDecoder(f).Decode(a)
	if err != nil {
		panic(err)
	}
}
