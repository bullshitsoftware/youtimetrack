package main

import (
	"fmt"

	"github.com/bullshitsoftware/youtimetrack/internal/app"
)

type App interface {
	Save() string
}

var a App = &app.App{}

func main() {
	p := a.Save()
	fmt.Println("Created", p, "edit it with your favorite text editor")
}
