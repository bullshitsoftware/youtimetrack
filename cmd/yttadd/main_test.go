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
	types := []youtrack.Type{
		{
			Id:   "123",
			Name: "Development",
		},
		{
			Id:   "456",
			Name: "DevOps",
		},
	}

	return types, nil
}

func (yt *YoutrackStub) Add(itemType youtrack.Type, issueId, duration, text string) error {
	if itemType.Id != "123" {
		panic("unexpected type")
	}

	if issueId != "XY-123" {
		panic("unexpected issue")
	}

	if duration != "1h" {
		panic("unexpected dureation")
	}

	if text != "did something" {
		panic("unexpected text")
	}

	return nil
}

func (yt *YoutrackStub) Delete(issueId, itemId string) error {
	panic("unexpected call")
}

func Example() {
	a = &AppStub{}

	os.Args = []string{"yttadd"}
	main()

	os.Args = []string{"yttadd", "deve", "XY-123", "1h", "did something"}
	main()

	os.Args = []string{"yttadd", "strange", "XY-123", "1h", "did something"}
	main()

	// Output:
	// Usage of yttadd type issue duration comment:
	//   - type, work type prefix (e.g., develop)
	//   - issue, issue number (e.g., XY-123)
	//   - duration, spent time in YouTrack format (e.g., 1h 30m)
	//   - comment, work item description (e.g., did something cool)
	// Time tracked
	// No work type found, available: Development, DevOps
}
