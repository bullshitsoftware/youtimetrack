package calendar

import "time"

type Calendar struct {
	DayDur  int `json:"day_dur"`
	SDayDur int `json:"short_day_dur"`

	Weekends  []time.Weekday `json:"weekends"`
	Workdays  []string       `json:"workdays"`
	SWorkdays []string       `json:"short_workdays"`
	Holidays  []string       `json:"holidays"`
}

type Period struct {
	Start time.Time
	End   time.Time
}

func (p *Period) ParseStart(s string) error {
	start, err := time.Parse("2006-01-02", s)
	if err != nil {
		return err
	}

	p.Start = start
	return nil
}

func (p *Period) ParseEnd(s string) error {
	end, err := time.Parse("2006-01-02", s)
	if err != nil {
		return err
	}

	p.End = time.Date(end.Year(), end.Month(), end.Day(), 23, 59, 59, 0, time.UTC)
	return nil
}

func (c *Calendar) Period(now time.Time) *Period {
	curYear, curMonth, _ := now.Date()
	start := time.Date(curYear, curMonth, 1, 0, 0, 0, 0, time.UTC)
	_, _, lastDay := start.AddDate(0, 1, -1).Date()
	end := time.Date(curYear, curMonth, lastDay, 23, 59, 59, 59, time.UTC)

	return &Period{start, end}
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
