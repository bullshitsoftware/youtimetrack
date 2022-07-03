package main

import (
	"os"
	"strings"
)

func main() {
	app := Default()

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
	case "init":
		Init(app)
	case "details":
		app.ReadConfig()
		Details(app, args)
	case "summary":
		app.ReadConfig()
		Summary(app, args)
	}
}
