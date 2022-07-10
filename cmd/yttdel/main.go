package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/bullshitsoftware/youtimetrack/internal/app"
)

type App interface {
	Load()
	NewYoutrack() app.Youtrack
}

var a App = &app.App{}

func main() {
	flag.Usage = func() {
		fmt.Printf("Usage of %s issue item_id:\n", os.Args[0])
		fmt.Println("  - issue, issue number (e.g., XY-123)")
		fmt.Println("  - item_id, work item id (e.g., 110-12312)")
	}
	flag.Parse()
	if flag.NArg() != 2 {
		flag.Usage()
		return
	}

	a.Load()
	yt := a.NewYoutrack()

	args := flag.Args()
	err := yt.Delete(args[0], args[1])
	app.ExitOnError(err)
	fmt.Println("Item deleted")
}
