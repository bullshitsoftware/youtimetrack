package app

import (
	"time"
)

type Config struct {
	Youtrack YoutrackConfig `json:"youtrack"`
	Calendar CalendarConfig `json:"calendar"`
}

type YoutrackConfig struct {
	BaseUrl string `json:"base_url"`
	Token   string `json:"token"`
	Author  string `json:"author"`
}

type CalendarConfig struct {
	DayDur  int `json:"day_dur"`
	SDayDur int `json:"short_day_dur"`

	Weekends  []time.Weekday     `json:"weekends"`
	Workdays  []string           `json:"workdays"`
	SWorkdays []string           `json:"short_workdays"`
	Holidays  []string           `json:"holidays"`
	Vacations []CalendarVacation `json:"vacations"`
}

type CalendarVacation struct {
	Start string `json:"start"`
	End   string `json:"end"`
}
