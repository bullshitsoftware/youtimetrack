package main

import (
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

type WorkItemDuration struct {
	Minutes int `json:"minutes"`
}

type WorkItem struct {
	Issue    Issue            `json:"issue"`
	Date     int64            `json:"date"`
	Duration WorkItemDuration `json:"duration"`
	Text     string           `json:"text"`
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
