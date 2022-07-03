package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"path"
	"strconv"
	"time"
)

const config = "config.json"

var (
	now  time.Time
	home string
)

func init() {
	now = time.Now()
	home = path.Join(os.Getenv("HOME"), ".config", "youtimetrack")
}

func main() {
	app := Default()

	if len(os.Args) == 2 && os.Args[1] == "init" {
		app.SaveConfig()
		fmt.Println("Created " + path.Join(home, config))

		return
	}

	app.ReadConfig()
	start, end := app.Calendar.Period(now)
	month := app.Calendar.Calc(start, end)
	today := app.Calendar.Calc(start, now)
	worked := app.Youtrack.Fetch(start, end)
	fmt.Printf("%s / %s / %s (worked / today / month)\n",
		FormatMinutes(worked),
		FormatMinutes(today),
		FormatMinutes(month),
	)
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

type App struct {
	Youtrack Youtrack `json:"youtrack"`
	Calendar Calendar `json:"calendar"`
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

type Youtrack struct {
	BaseUrl string `json:"base_url"`
	Token   string `json:"token"`
	Author  string `json:"author"`
}

func (t *Youtrack) Fetch(start, end time.Time) int {
	q := url.Values{}
	q.Add("fields", "duration(minutes)")
	q.Add("author", t.Author)
	q.Add("start", strconv.FormatInt(start.UnixMilli(), 10))
	q.Add("end", strconv.FormatInt(end.UnixMilli(), 10))

	u, err := url.Parse(t.BaseUrl + "/workItems")
	if err != nil {
		panic(err)
	}
	u.RawQuery = q.Encode()

	req, err := http.NewRequest("GET", u.String(), nil)
	if err != nil {
		panic(err)
	}
	req.Header.Set("Authorization", "Bearer "+t.Token)

	client := http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	type WorkItemDuration struct {
		Minutes int `json:"minutes"`
	}

	type WorkItem struct {
		Duration WorkItemDuration `json:"duration"`
	}

	items := []WorkItem{}
	err = json.NewDecoder(resp.Body).Decode(&items)
	if err != nil {
		panic(err)
	}

	m := 0
	for _, i := range items {
		m += i.Duration.Minutes
	}

	return m
}

type Calendar struct {
	DayDur  int `json:"day_dur"`
	SDayDur int `json:"short_day_dur"`

	Weekends  []time.Weekday `json:"weekends"`
	Workdays  []string       `json:"workdays"`
	SWorkdays []string       `json:"short_workdays"`
	Holidays  []string       `json:"holidays"`
}

func (c *Calendar) Period(now time.Time) (time.Time, time.Time) {
	curYear, curMonth, _ := now.Date()
	start := time.Date(curYear, curMonth, 1, 0, 0, 0, 0, time.UTC)
	_, _, lastDay := start.AddDate(0, 1, -1).Date()
	end := time.Date(curYear, curMonth, lastDay, 23, 59, 59, 59, time.UTC)

	return start, end
}

func (c *Calendar) Calc(start, end time.Time) int {
	m := 0
OUTER:
	for cur := start; cur.Unix() <= end.Unix(); cur = cur.AddDate(0, 0, 1) {
		d := cur.Format("2006-01-02")
		for _, holiday := range c.Holidays {
			if holiday == d {
				continue OUTER
			}
		}

		for _, workday := range c.Workdays {
			if workday == d {
				m += c.DayDur
				continue OUTER
			}
		}

		for _, workday := range c.SWorkdays {
			if workday == d {
				m += c.SDayDur
				continue OUTER
			}
		}

		for _, weekend := range c.Weekends {
			if weekend == cur.Weekday() {
				continue OUTER
			}
		}

		m += c.DayDur
	}

	return m
}

func FormatMinutes(m int) string {
	s := fmt.Sprintf("%dh", m/60)
	if m%60 > 0 {
		s += fmt.Sprintf(" %dm", m%60)
	}

	return s
}
