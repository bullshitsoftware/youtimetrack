package main

import (
	"flag"
	"fmt"
	"os"
	"strings"

	"github.com/bullshitsoftware/youtimetrack/internal/app"
)

type App interface {
	Load()
	NewYoutrack() app.Youtrack
}

var a App = &app.App{}

func main() {
	flag.Usage = func() {
		fmt.Printf("Usage of %s type issue duration comment:\n", os.Args[0])
		fmt.Println("  - type, work type prefix (e.g., develop)")
		fmt.Println("  - issue, issue number (e.g., XY-123)")
		fmt.Println("  - duration, spent time in YouTrack format (e.g., 1h 30m)")
		fmt.Println("  - comment, work item description (e.g., did something cool)")
	}
	flag.Parse()
	if flag.NArg() != 4 {
		flag.Usage()
		return
	}

	a.Load()
	yt := a.NewYoutrack()

	args := flag.Args()
	typeName := strings.ToLower(args[0])
	issue := args[1]
	duration := args[2]
	text := args[3]

	types, err := yt.WorkItemTypes()
	app.ExitOnError(err)
	aTypes := []string{}
	for _, i := range types {
		s := strings.ToLower(i.Name)
		if strings.HasPrefix(s, typeName) {
			err = yt.Add(i, issue, duration, text)
			app.ExitOnError(err)
			fmt.Println("Time tracked")
			return
		}
		aTypes = append(aTypes, i.Name)
	}
	fmt.Println("No work type found, available:", strings.Join(aTypes, ", "))
}
