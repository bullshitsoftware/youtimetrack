package youtrack

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/url"
	"strconv"
	"time"
)

type Client struct {
	BaseUrl string `json:"base_url"`
	Token   string `json:"token"`
	Author  string `json:"author"`
}

type Issue struct {
	IdReadable string `json:"idReadable"`
	Summary    string `josn:"summary"`
}

type Duration struct {
	Minutes      int    `json:"minutes,omitempty"`
	Presentation string `json:"presentation,omitempty"`
}

type Type struct {
	Id   string `json:"id"`
	Name string `json:"name,omitempty"`
}

type WorkItem struct {
	Issue    *Issue   `json:"issue,omitempty"`
	Date     int64    `json:"date,omitempty"`
	Duration Duration `json:"duration"`
	Type     Type     `json:"type"`
	Text     string   `json:"text"`
}

func (c *Client) Fetch(start, end time.Time) []WorkItem {
	q := url.Values{}
	q.Add("fields", "issue(idReadable,summary),date,duration(minutes),text")
	q.Add("author", c.Author)
	q.Add("start", strconv.FormatInt(start.UnixMilli(), 10))
	q.Add("end", strconv.FormatInt(end.UnixMilli(), 10))

	u, err := url.Parse(c.BaseUrl + "/workItems")
	if err != nil {
		panic(err)
	}
	u.RawQuery = q.Encode()

	req, err := http.NewRequest("GET", u.String(), nil)
	if err != nil {
		panic(err)
	}
	req.Header.Set("Authorization", "Bearer "+c.Token)

	client := http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	items := []WorkItem{}
	err = json.NewDecoder(resp.Body).Decode(&items)
	if err != nil {
		panic(err)
	}

	return items
}

func (c *Client) WorkItemTypes() []Type {
	q := url.Values{}
	q.Add("fields", "id,name")

	u, err := url.Parse(c.BaseUrl + "/admin/timeTrackingSettings/workItemTypes")
	if err != nil {
		panic(err)
	}
	u.RawQuery = q.Encode()

	req, err := http.NewRequest("GET", u.String(), nil)
	if err != nil {
		panic(err)
	}
	req.Header.Set("Authorization", "Bearer "+c.Token)

	client := http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	items := []Type{}
	err = json.NewDecoder(resp.Body).Decode(&items)
	if err != nil {
		panic(err)
	}

	return items
}

func (c *Client) Add(itemType Type, issueId, duration, text string) {
	body := WorkItem{
		Type:     itemType,
		Duration: Duration{Presentation: "1h"},
		Text:     text,
	}
	b, _ := json.Marshal(body)

	req, err := http.NewRequest("POST", c.BaseUrl+"/issues/"+issueId+"/timeTracking/workItems", bytes.NewReader(b))
	if err != nil {
		panic(err)
	}
	req.Header.Set("Authorization", "Bearer "+c.Token)
	req.Header.Set("Content-Type", "application/json")

	client := http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()
}
