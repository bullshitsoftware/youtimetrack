package app

import (
	"io"
	"os"
	"path"
	"testing"
	"time"

	"github.com/bullshitsoftware/youtimetrack/internal/pkg/calendar"
	"github.com/bullshitsoftware/youtimetrack/internal/pkg/youtrack"
	"github.com/stretchr/testify/assert"
)

func TestApp_Load(t *testing.T) {
	oldHome := home
	home = path.Join(os.TempDir(), "ytt")
	assert := assert.New(t)

	src, err := os.Open("config_stub.json")
	assert.NoError(err)
	defer src.Close()

	err = os.MkdirAll(home, 0700)
	assert.NoError(err)
	dest, err := os.Create(path.Join(home, config))
	assert.NoError(err)
	defer dest.Close()

	_, err = io.Copy(dest, src)
	assert.NoError(err)

	app := App{}
	app.Load()

	yt := app.cfg.Youtrack
	assert.Equal("https://localhost/api", yt.BaseUrl)
	assert.Equal("token", yt.Token)
	assert.Equal("author_id", yt.Author)
	cal := app.cfg.Calendar
	assert.Equal(480, cal.DayDur)
	assert.Equal(420, cal.SDayDur)
	assert.Equal([]time.Weekday{time.Sunday, time.Saturday}, cal.Weekends)
	assert.Equal([]string{"2022-01-01"}, cal.Workdays)
	assert.Equal([]string{"2022-01-02"}, cal.SWorkdays)
	assert.Equal([]string{"2022-01-03"}, cal.Holidays)

	os.RemoveAll(home)
	home = oldHome
}

func TestApp_Save(t *testing.T) {
	oldHome := home
	home = path.Join(os.TempDir(), "ytt")
	assert := assert.New(t)

	app := App{
		cfg: Config{
			Youtrack: YoutrackConfig{
				BaseUrl: "https://localhost/api",
				Token:   "token",
				Author:  "author_id",
			},
			Calendar: CalendarConfig{
				DayDur:    480,
				SDayDur:   420,
				Weekends:  []time.Weekday{time.Sunday, time.Saturday},
				Workdays:  []string{"2022-01-01"},
				SWorkdays: []string{"2022-01-02"},
				Holidays:  []string{"2022-01-03"},
				Vacations: []CalendarVacation{{"2022-01-04", "2022-01-06"}},
			},
		},
	}
	app.Save()

	stub, err := os.ReadFile("config_stub.json")
	assert.NoError(err)
	saved, err := os.ReadFile(path.Join(home, config))
	assert.NoError(err)
	assert.Equal(string(stub), string(saved))

	os.RemoveAll(home)
	home = oldHome
}

func TestApp_NewCalendar(t *testing.T) {
	assert := assert.New(t)

	app := App{}
	app.cfg.Calendar.DayDur = 8
	app.cfg.Calendar.SDayDur = 7
	app.cfg.Calendar.Weekends = []time.Weekday{time.Sunday}
	app.cfg.Calendar.Workdays = []string{"2007-01-02"}
	app.cfg.Calendar.SWorkdays = []string{"2007-01-03"}
	app.cfg.Calendar.Holidays = []string{"2007-01-04"}
	app.cfg.Calendar.Vacations = []CalendarVacation{{"2022-01-05", "2022-01-07"}}

	var cal Calendar
	var err error
	cal, err = app.NewCalendar()
	assert.NoError(err)
	c := cal.(*calendar.Calendar)
	assert.Equal(8, c.DayDur)
	assert.Equal(7, c.SDayDur)
	assert.Equal([]time.Weekday{time.Sunday}, c.Weekends)
	assert.Equal([]string{"2007-01-02"}, c.Workdays)
	assert.Equal([]string{"2007-01-03"}, c.SWorkdays)
	assert.Equal([]string{"2022-01-05", "2022-01-06", "2022-01-07", "2007-01-04"}, c.Holidays)

	app.cfg.Calendar.Vacations = []CalendarVacation{{"start", "2022-01-07"}}
	cal, err = app.NewCalendar()
	assert.Error(err)
	assert.ErrorContains(err, "cannot parse \"start\"")
	assert.Nil(cal)

	app.cfg.Calendar.Vacations = []CalendarVacation{{"2022-01-05", "end"}}
	cal, err = app.NewCalendar()
	assert.Error(err)
	assert.ErrorContains(err, "cannot parse \"end\"")
	assert.Nil(cal)
}

func TestApp_NewYoutrack(t *testing.T) {
	assert := assert.New(t)

	app := App{}
	app.cfg.Youtrack.BaseUrl = "https://localhost"
	app.cfg.Youtrack.Token = "token"
	app.cfg.Youtrack.Author = "author"

	yt := app.NewYoutrack().(*youtrack.Client)
	assert.Equal("https://localhost", yt.BaseUrl)
	assert.Equal("token", yt.Token)
	assert.Equal("author", yt.Author)
}
