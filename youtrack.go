package main

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/url"
	"strconv"
	"time"
)

type Youtrack struct {
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

func (t *Youtrack) Fetch(start, end time.Time) []WorkItem {
	q := url.Values{}
	q.Add("fields", "issue(idReadable,summary),date,duration(minutes),text")
	q.Add("author", t.Author)
	q.Add("start", strconv.FormatInt(start.UnixMilli(), 10))
	q.Add("end", strconv.FormatInt(end.UnixMilli(), 10))

	u, err := url.Parse(t.BaseUrl + "/workItems")
	if err != nil {
		panic(err)
	}
	u.RawQuery = q.Encode()

	req, err := http.NewRequest("GET", u.String(), nil)
	if err != nil {
		panic(err)
	}
	req.Header.Set("Authorization", "Bearer "+t.Token)

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

func (t *Youtrack) WorkItemTypes() []Type {
	q := url.Values{}
	q.Add("fields", "id,name")

	u, err := url.Parse(t.BaseUrl + "/admin/timeTrackingSettings/workItemTypes")
	if err != nil {
		panic(err)
	}
	u.RawQuery = q.Encode()

	req, err := http.NewRequest("GET", u.String(), nil)
	if err != nil {
		panic(err)
	}
	req.Header.Set("Authorization", "Bearer "+t.Token)

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

func (t *Youtrack) Add(itemType Type, issueId, duration, text string) {
	body := WorkItem{
		Type:     itemType,
		Duration: Duration{Presentation: "1h"},
		Text:     text,
	}
	b, _ := json.Marshal(body)

	req, err := http.NewRequest("POST", t.BaseUrl+"/issues/"+issueId+"/timeTracking/workItems", bytes.NewReader(b))
	if err != nil {
		panic(err)
	}
	req.Header.Set("Authorization", "Bearer "+t.Token)
	req.Header.Set("Content-Type", "application/json")

	client := http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()
}
