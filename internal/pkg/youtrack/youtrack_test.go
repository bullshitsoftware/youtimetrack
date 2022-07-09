package youtrack

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestYoutrack_WorkItemTypes(t *testing.T) {
	assert := assert.New(t)
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(http.MethodGet, r.Method)
		assert.Equal("/admin/timeTrackingSettings/workItemTypes", r.URL.Path)
		assert.Len(r.Header["Authorization"], 1)
		assert.Equal("Bearer token", r.Header["Authorization"][0])

		body := []Type{
			{Id: "1", Name: "name"},
		}
		json.NewEncoder(w).Encode(body)
	}))
	defer ts.Close()

	yt := Client{
		BaseUrl: ts.URL,
		Token:   "token",
		Author:  "id",
	}
	types := yt.WorkItemTypes()
	assert.Len(types, 1)
	assert.Equal("1", types[0].Id)
	assert.Equal("name", types[0].Name)
}

func TestYoutrack_Add(t *testing.T) {
	assert := assert.New(t)
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(http.MethodPost, r.Method)
		assert.Equal("/issues/XY-123/timeTracking/workItems", r.URL.Path)
		assert.Len(r.Header["Authorization"], 1)
		assert.Equal("Bearer token", r.Header["Authorization"][0])
		assert.Len(r.Header["Content-Type"], 1)
		assert.Equal("application/json", r.Header["Content-Type"][0])

		item := WorkItem{}
		err := json.NewDecoder(r.Body).Decode(&item)
		assert.NoError(err)
		assert.Equal("1", item.Type.Id)
		assert.Equal("1h", item.Duration.Presentation)
		assert.Equal("text", item.Text)
	}))
	defer ts.Close()

	yt := Client{
		BaseUrl: ts.URL,
		Token:   "token",
		Author:  "id",
	}
	yt.Add(Type{Id: "1"}, "XY-123", "1h", "text")
}
