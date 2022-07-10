package app

import (
	"encoding/json"
	"fmt"
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

func (a *App) SaveConfig() (string, error) {
	err := os.MkdirAll(home, 0700)
	if err != nil {
		return "", err
	}

	p := path.Join(home, config)
	f, err := os.Create(p)
	if err != nil {
		return "", err
	}
	defer f.Close()

	enc := json.NewEncoder(f)
	enc.SetIndent("", "  ")
	err = enc.Encode(a)
	if err != nil {
		return "", err
	}

	return p, nil
}

func (a *App) ReadConfig() error {
	f, err := os.Open(path.Join(home, config))
	if err != nil {
		return err
	}
	defer f.Close()

	err = json.NewDecoder(f).Decode(a)
	if err != nil {
		return err
	}

	return nil
}

func ExitOnError(err error) {
	if err != nil {
		printError(err)
		os.Exit(1)
	}
}

func printError(err error) {
	fmt.Println("Error:", err)
}

func FormatMinutes(m int) string {
	s := fmt.Sprintf("%dh", m/60)
	if m%60 > 0 {
		s += fmt.Sprintf("%dm", m%60)
	}

	return s
}
