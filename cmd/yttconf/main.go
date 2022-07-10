package main

import (
	"fmt"

	"github.com/bullshitsoftware/youtimetrack/internal/app"
)

func main() {
	a := app.Default()
	p, err := a.SaveConfig()
	app.ExitOnError(err)
	fmt.Println("Created", p, "edit it with your favorite text editor")
}
