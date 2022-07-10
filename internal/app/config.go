package app

import (
	"encoding/json"
	"io"
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

	Weekends  []time.Weekday `json:"weekends"`
	Workdays  []string       `json:"workdays"`
	SWorkdays []string       `json:"short_workdays"`
	Holidays  []string       `json:"holidays"`
}

func (c *Config) LoadJson(r io.Reader) error {
	err := json.NewDecoder(r).Decode(c)
	if err != nil {
		return err
	}
	return nil
}

func (c *Config) SaveJson(w io.Writer) error {
	enc := json.NewEncoder(w)
	enc.SetIndent("", "  ")
	err := enc.Encode(c)
	if err != nil {
		return err
	}

	return nil
}
