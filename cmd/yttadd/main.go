package main

import (
	"errors"
	"os"
	"strings"

	"github.com/bullshitsoftware/youtimetrack/internal/app"
	yt "github.com/bullshitsoftware/youtimetrack/internal/pkg/youtrack"
)

func main() {
	if len(os.Args) != 4 {
		app.ExitOnError(errors.New("invalid arguments number"))
	}

	a := app.Default()
	err := a.ReadConfig()
	app.ExitOnError(err)
	typeName := strings.ToLower(os.Args[0])
	types, err := a.Youtrack.WorkItemTypes()
	app.ExitOnError(err)
	var t yt.Type
	for _, i := range types {
		s := strings.ToLower(i.Name)
		if strings.HasPrefix(s, typeName) {
			t = i
			break
		}
	}
	issue := os.Args[1]
	duration := os.Args[2]
	text := os.Args[3]

	a.Youtrack.Add(t, issue, duration, text)
}
