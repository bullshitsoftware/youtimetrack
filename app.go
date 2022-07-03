package main

import (
	"encoding/json"
	"os"
	"path"
	"time"
)

const config = "config.json"

var (
	now  time.Time
	home string
)

func init() {
	now = time.Now()
	home = path.Join(os.Getenv("HOME"), ".config", "ytt")
}

type App struct {
	Youtrack Youtrack `json:"youtrack"`
	Calendar Calendar `json:"calendar"`
}

func Default() *App {
	return &App{
		Youtrack{
			BaseUrl: "http://localhost:2378/api",
			Token:   "your-token",
			Author:  "your-user-uuid",
		},
		Calendar{
			DayDur:    8 * 60,
			SDayDur:   7 * 60,
			Weekends:  []time.Weekday{time.Saturday, time.Sunday},
			Workdays:  []string{"2022-02-05"},
			SWorkdays: []string{"2022-02-21"},
			Holidays:  []string{"2022-02-22"},
		},
	}
}

func (a *App) SaveConfig() {
	err := os.MkdirAll(home, 0700)
	if err != nil {
		panic(err)
	}

	f, err := os.Create(path.Join(home, config))
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
}

func (a *App) ReadConfig() {
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
