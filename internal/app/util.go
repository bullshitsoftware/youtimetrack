package app

import (
	"fmt"
	"os"

	yt "github.com/bullshitsoftware/youtimetrack/internal/pkg/youtrack"
)

var exit = os.Exit

func ExitOnError(err error) {
	if err != nil {
		PrintError(err)
		exit(1)
	}
}

func PrintError(err error) {
	fmt.Println("Error:", err)
	if v, ok := err.(*yt.UnexpectedResponseError); ok {
		fmt.Println("response body:", string(v.Body))
	}
}

func FormatMinutes(m int) string {
	s := fmt.Sprintf("%dh", m/60)
	if m%60 > 0 {
		s += fmt.Sprintf("%dm", m%60)
	}

	return s
}
