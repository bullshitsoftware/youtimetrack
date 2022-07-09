package calendar

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestPeriod_ParseStart(t *testing.T) {
	assert := assert.New(t)
	p := Period{}
	var err error
	err = p.ParseStart("2022-01-02")
	assert.NoError(err)
	assert.Equal(time.Date(2022, 1, 2, 0, 0, 0, 0, time.UTC), p.Start)

	err = p.ParseStart("01-02-2022")
	assert.Error(err)
	assert.ErrorContains(err, "cannot parse")
}

func TestPeriod_ParseEnd(t *testing.T) {
	assert := assert.New(t)
	p := Period{}
	var err error
	err = p.ParseEnd("2022-01-02")
	assert.NoError(err)
	assert.Equal(time.Date(2022, 1, 2, 23, 59, 59, 0, time.UTC), p.End)

	err = p.ParseEnd("01-02-2022")
	assert.Error(err)
	assert.ErrorContains(err, "cannot parse")
}

func TestCalendar_Period(t *testing.T) {
	assert := assert.New(t)
	c := Calendar{}
	p := c.Period(time.Date(2022, 2, 1, 3, 4, 5, 6, time.UTC))
	assert.Equal(time.Date(2022, 2, 1, 0, 0, 0, 0, time.UTC), p.Start)
	assert.Equal(time.Date(2022, 2, 28, 23, 59, 59, 0, time.UTC), p.End)
}

func TestCalendar_Calc(t *testing.T) {
	assert := assert.New(t)
	c := Calendar{
		DayDur:    8 * 60,
		SDayDur:   7 * 60,
		Holidays:  []string{"2022-02-01"},
		Workdays:  []string{"2022-02-05"},
		SWorkdays: []string{"2022-02-04"},
		Weekends:  []time.Weekday{time.Sunday, time.Monday},
	}
	m := c.Calc(time.Date(2022, 2, 1, 0, 0, 0, 0, time.UTC), time.Date(2022, 2, 6, 23, 59, 59, 0, time.UTC))
	// 6 days total = 6*8*60 = 2 880 total minutes
	// two weekends, but one is remapped as a workday -8*60=-480 minutes
	// one half-holiday -60 minutes
	// and one public holiday -8*60=-480 minutes
	// 2 880 -480 -60 -480 = 1860
	assert.Equal(1860, m)
}
