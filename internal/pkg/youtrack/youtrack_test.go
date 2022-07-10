package youtrack

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestNewUnexpectedResponseError(t *testing.T) {
	assert := assert.New(t)
	r := http.Response{
		Status: "status",
		Body:   ioutil.NopCloser(strings.NewReader("body")),
	}
	err := NewUnexpectedResponseError(&r)
	assert.Error(err)
	assert.Equal("status", err.Status)
	assert.Equal("body", string(err.Body))
}

func TestUnexpectedResponseError_Error(t *testing.T) {
	err := UnexpectedResponseError{
		Status: "status",
		Body:   []byte("body"),
	}
	assert.Equal(t, "unexpected response status: status", err.Error())
}

func TestClient_WorkItems(t *testing.T) {
	assert := assert.New(t)

	start := time.Date(2022, 2, 1, 0, 0, 0, 0, time.UTC)
	end := time.Date(2022, 2, 28, 23, 59, 59, 0, time.UTC)
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(http.MethodGet, r.Method)
		assert.Equal("/workItems", r.URL.Path)

		q := r.URL.Query()
		assert.Len(q, 4)
		assert.Equal([]string{"issue(idReadable,summary),date,duration(minutes),text"}, q["fields"])
		assert.Equal([]string{"id"}, q["author"])
		assert.Equal([]string{strconv.FormatInt(start.UnixMilli(), 10)}, q["start"])
		assert.Equal([]string{strconv.FormatInt(end.UnixMilli(), 10)}, q["end"])

		items := []WorkItem{
			{
				Issue:    &Issue{IdReadable: "XY-123", Summary: "Issue summary"},
				Date:     123,
				Duration: Duration{Minutes: 80},
				Text:     "Text1\nText2",
			},
		}
		json.NewEncoder(w).Encode(items)
	}))
	defer ts.Close()

	c := Client{
		BaseUrl: ts.URL,
		Token:   "token",
		Author:  "id",
	}
	var items []WorkItem
	var err error
	items, err = c.WorkItems(start, end)
	assert.NoError(err)
	assert.Len(items, 1)
	assert.NotNil(items[0].Issue)
	assert.Equal("XY-123", items[0].Issue.IdReadable)
	assert.Equal("Issue summary", items[0].Issue.Summary)
	assert.Equal(int64(123), items[0].Date)
	assert.Equal(80, items[0].Duration.Minutes)
	assert.Equal("Text1\nText2", items[0].Text)

	ts.Config.Handler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "not json")
	})
	items, err = c.WorkItems(time.Now(), time.Now())
	assert.Error(err)
	assert.ErrorContains(err, "invalid character")
	assert.Len(items, 0)

	ts.Config.Handler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusCreated)
		fmt.Fprintln(w, "[]")
	})
	items, err = c.WorkItems(time.Now(), time.Now())
	assert.Error(err)
	assert.IsType(&UnexpectedResponseError{}, err)
	assert.Equal("[]\n", string(err.(*UnexpectedResponseError).Body))
	assert.ErrorContains(err, "unexpected response status")
	assert.Len(items, 0)

	c.BaseUrl = "ptth://localhost"
	items, err = c.WorkItems(time.Now(), time.Now())
	assert.Error(err)
	assert.ErrorContains(err, "unsupported protocol scheme")
	assert.Len(items, 0)
}

func TestClient_WorkItemTypes(t *testing.T) {
	assert := assert.New(t)
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(http.MethodGet, r.Method)
		assert.Equal("/admin/timeTrackingSettings/workItemTypes", r.URL.Path)

		assert.Equal([]string{"id,name"}, r.URL.Query()["fields"])

		body := []Type{
			{Id: "1", Name: "name"},
		}
		json.NewEncoder(w).Encode(body)
	}))
	defer ts.Close()

	c := Client{
		BaseUrl: ts.URL,
		Token:   "token",
		Author:  "id",
	}
	var items []Type
	var err error
	items, err = c.WorkItemTypes()
	assert.NoError(err)
	assert.Len(items, 1)
	assert.Equal("1", items[0].Id)
	assert.Equal("name", items[0].Name)

	ts.Config.Handler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "not json")
	})
	items, err = c.WorkItemTypes()
	assert.Error(err)
	assert.ErrorContains(err, "invalid character")
	assert.Len(items, 0)

	ts.Config.Handler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusCreated)
		fmt.Fprintln(w, "[]")
	})
	items, err = c.WorkItemTypes()
	assert.Error(err)
	assert.IsType(&UnexpectedResponseError{}, err)
	assert.Equal("[]\n", string(err.(*UnexpectedResponseError).Body))
	assert.ErrorContains(err, "unexpected response status")
	assert.Len(items, 0)

	c.BaseUrl = "ptth://localhost"
	items, err = c.WorkItemTypes()
	assert.Error(err)
	assert.ErrorContains(err, "unsupported protocol scheme")
	assert.Len(items, 0)
}

func TestClient_Add(t *testing.T) {
	assert := assert.New(t)
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(http.MethodPost, r.Method)
		assert.Equal("/issues/XY-123/timeTracking/workItems", r.URL.Path)

		item := WorkItem{}
		err := json.NewDecoder(r.Body).Decode(&item)
		assert.NoError(err)
		assert.Equal("1", item.Type.Id)
		assert.Equal("1h", item.Duration.Presentation)
		assert.Equal("text", item.Text)
	}))
	defer ts.Close()

	c := Client{
		BaseUrl: ts.URL,
		Token:   "token",
		Author:  "id",
	}
	var err error
	err = c.Add(Type{Id: "1"}, "XY-123", "1h", "text")
	assert.NoError(err)

	c.BaseUrl = "ptth://localhost"
	err = c.Add(Type{Id: "1"}, "XY-123", "1h", "text")
	assert.Error(err)
	assert.ErrorContains(err, "unsupported protocol scheme")
}

func TestClient_Delete(t *testing.T) {
	assert := assert.New(t)
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(http.MethodDelete, r.Method)
		assert.Equal("/issues/XY-123/timeTracking/workItems/321", r.URL.Path)
	}))
	defer ts.Close()

	c := Client{
		BaseUrl: ts.URL,
		Token:   "token",
		Author:  "id",
	}
	var err error
	err = c.Delete("XY-123", "321")
	assert.NoError(err)

	ts.Config.Handler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusCreated)
	})
	err = c.Delete("XY-123", "321")
	assert.Error(err)
	assert.IsType(&UnexpectedResponseError{}, err)
	assert.Equal("", string(err.(*UnexpectedResponseError).Body))

	c.BaseUrl = "ptth://localhost"
	err = c.Delete("XY-123", "321")
	assert.Error(err)
	assert.ErrorContains(err, "unsupported protocol scheme")
}

func TestClient_request(t *testing.T) {
	assert := assert.New(t)
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(http.MethodPost, r.Method)
		assert.Equal("/test", r.URL.Path)

		assert.Equal([]string{"Bearer token"}, r.Header["Authorization"])
		assert.Equal([]string{"application/json"}, r.Header["Content-Type"])

		q := r.URL.Query()
		assert.Len(q["foo"], 1)
		assert.Equal("bar", q["foo"][0])

		b, _ := io.ReadAll(r.Body)
		assert.Equal("test", string(b))
	}))
	c := Client{
		BaseUrl: ts.URL,
		Token:   "token",
		Author:  "id",
	}

	var resp *http.Response
	var err error

	body := bytes.NewReader([]byte("test"))
	resp, err = c.request(http.MethodPost, "/test", url.Values{"foo": []string{"bar"}}, body)
	assert.NoError(err)
	assert.NotNil(resp)

	resp, err = c.request("[", "/test", nil, nil)
	assert.Nil(resp)
	assert.Error(err)
	assert.ErrorContains(err, "invalid method")

	c.BaseUrl = "http://local host"
	resp, err = c.request(http.MethodGet, "/t", nil, nil)
	assert.Nil(resp)
	assert.Error(err)
	assert.ErrorContains(err, "invalid character")

	c.BaseUrl = "ptth://localhost"
	resp, err = c.request(http.MethodGet, "/t", nil, nil)
	assert.Nil(resp)
	assert.Error(err)
	assert.ErrorContains(err, "unsupported protocol scheme")
}
