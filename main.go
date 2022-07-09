package main

import (
	"os"
	"strings"

	"github.com/bullshitsoftware/youtimetrack/internal/app"
)

func main() {
	app := app.Default()

	var args []string
	var cmd string
	if len(os.Args) > 1 && !strings.HasPrefix(os.Args[1], "-") {
		cmd = os.Args[1]
		args = os.Args[2:]
	} else {
		cmd = "summary"
		args = os.Args[1:]
	}

	switch cmd {
	case "i", "init":
		Init(app)
	case "d", "details":
		app.ReadConfig(home)
		Details(app, args)
	case "s", "summary":
		app.ReadConfig(home)
		Summary(app, args)
	case "a", "add":
		app.ReadConfig(home)
		Add(app, args)
	}
}
