package main

import (
	"os"
	"time"

	"github.com/bullshitsoftware/youtimetrack/internal/app"
	"github.com/bullshitsoftware/youtimetrack/internal/pkg/youtrack"
)

type AppStub struct{}

func (s *AppStub) Load() {}

func (s *AppStub) NewYoutrack() app.Youtrack {
	return &YoutrackStub{}
}

type YoutrackStub struct{}

func (yt *YoutrackStub) WorkItems(start, end time.Time) ([]youtrack.WorkItem, error) {
	panic("unexpected call")
}

func (yt *YoutrackStub) WorkItemTypes() ([]youtrack.Type, error) {
	panic("unexpected call")
}

func (yt *YoutrackStub) Add(itemType youtrack.Type, issueId, duration, text string) error {
	panic("unexpected call")
}

func (yt *YoutrackStub) Delete(issueId, itemId string) error {
	if issueId != "XY-123" {
		panic("unexpected issue id " + issueId)
	}

	if itemId != "110-12312" {
		panic("unexpected item id " + itemId)
	}

	return nil
}

func Example() {
	a = &AppStub{}

	os.Args = []string{"yttdel"}
	main()

	os.Args = []string{"yttdel", "XY-123", "110-12312"}
	main()

	// Output:
	// Usage of yttdel issue item_id:
	//   - issue, issue number (e.g., XY-123)
	//   - item_id, work item id (e.g., 110-12312)
	// Item deleted
}
