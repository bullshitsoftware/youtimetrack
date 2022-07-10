package app

import (
	"errors"
	"testing"

	yt "github.com/bullshitsoftware/youtimetrack/internal/pkg/youtrack"
	"github.com/stretchr/testify/assert"
)

func TestExitOnError(t *testing.T) {
	ExitOnError(nil)
}

func ExamplePrintError() {
	PrintError(errors.New("test"))

	err := &yt.UnexpectedResponseError{
		Status: "status",
		Body:   []byte("body"),
	}
	PrintError(err)

	// Output:
	// Error: test
	// Error: unexpected response status: status
	// response body: body
}

func TestFormatMinutes(t *testing.T) {
	assert := assert.New(t)
	assert.Equal("0h35m", FormatMinutes(35))
	assert.Equal("1h40m", FormatMinutes(100))
	assert.Equal("2h", FormatMinutes(120))
}
