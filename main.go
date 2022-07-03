package main

import (
	"os"
)

func main() {
	app := Default()

	cmd := "summary"
	if len(os.Args) > 1 {
		cmd = os.Args[1]
	}
	switch cmd {
	case "init":
		Init(app)
	case "details":
		app.ReadConfig()
		Details(app)
	case "summary":
		app.ReadConfig()
		Summary(app)
	}
}
