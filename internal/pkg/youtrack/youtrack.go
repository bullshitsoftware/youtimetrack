package youtrack

import (
	"bytes"
	"encoding/json"
	"io"
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

type UnexpectedResponseError struct {
	Status string
	Body   []byte
}

func NewUnexpectedResponseError(r *http.Response) *UnexpectedResponseError {
	body, _ := io.ReadAll(r.Body)
	return &UnexpectedResponseError{
		Status: r.Status,
		Body:   body,
	}
}

func (e *UnexpectedResponseError) Error() string {
	return "unexpected response status: " + e.Status
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

func (c *Client) WorkItems(start, end time.Time) ([]WorkItem, error) {
	q := url.Values{}
	q.Add("fields", "issue(idReadable,summary),date,duration(minutes),text")
	q.Add("author", c.Author)
	q.Add("start", strconv.FormatInt(start.UnixMilli(), 10))
	q.Add("end", strconv.FormatInt(end.UnixMilli(), 10))

	resp, err := c.request(http.MethodGet, "/workItems", q, nil)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, NewUnexpectedResponseError(resp)
	}

	items := []WorkItem{}
	err = json.NewDecoder(resp.Body).Decode(&items)
	if err != nil {
		return nil, err
	}

	return items, nil
}

func (c *Client) WorkItemTypes() ([]Type, error) {
	q := url.Values{}
	q.Add("fields", "id,name")

	resp, err := c.request(http.MethodGet, "/admin/timeTrackingSettings/workItemTypes", q, nil)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, NewUnexpectedResponseError(resp)
	}

	items := []Type{}
	err = json.NewDecoder(resp.Body).Decode(&items)
	if err != nil {
		return nil, err
	}

	return items, nil
}

func (c *Client) Add(itemType Type, issueId, duration, text string) error {
	body := WorkItem{
		Type:     itemType,
		Duration: Duration{Presentation: "1h"},
		Text:     text,
	}
	b, _ := json.Marshal(body)

	resp, err := c.request(http.MethodPost, "/issues/"+issueId+"/timeTracking/workItems", nil, bytes.NewReader(b))
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	return nil
}

func (c *Client) request(method, path string, values url.Values, body io.Reader) (*http.Response, error) {
	u, err := url.Parse(c.BaseUrl + path)
	if err != nil {
		return nil, err
	}
	u.RawQuery = values.Encode()

	req, err := http.NewRequest(method, u.String(), body)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Authorization", "Bearer "+c.Token)
	req.Header.Set("Content-Type", "application/json")

	client := http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	return resp, nil
}
