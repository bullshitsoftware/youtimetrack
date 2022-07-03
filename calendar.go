package main

import "time"

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
