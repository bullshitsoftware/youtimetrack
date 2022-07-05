package main

import (
	"context"
	"encoding/json"
	"net/http"
	"os"
	"path"
	"strconv"
	"testing"
	"time"
)

var ytServer *http.Server

func init() {
	handler := http.NewServeMux()
	handler.HandleFunc("/api/workItems", func(w http.ResponseWriter, r *http.Request) {
		if v, ok := r.Header["Authorization"]; !ok && v[0] != "Bearer your-token" {
			http.Error(w, "Missing or invalid \"Authorization\" header", http.StatusBadRequest)

			return
		}

		if r.Method != "GET" {
			http.NotFound(w, r)

			return
		}

		q := r.URL.Query()
		if v, ok := q["fields"]; !ok || v[0] != "issue(idReadable,summary),date,duration(minutes)" {
			http.Error(w, "Missing or invalid \"fields\" query parameter", http.StatusBadRequest)

			return
		}

		if v, ok := q["author"]; !ok || v[0] != "your-user-uuid" {
			http.Error(w, "Missing or invalid \"author\" query parameter", http.StatusBadRequest)

			return
		}

		if _, ok := q["start"]; !ok {
			http.Error(w, "Missing \"start\" query parameter", http.StatusBadRequest)

			return
		}

		if _, ok := q["end"]; !ok {
			http.Error(w, "Missing \"end\" query parameter", http.StatusBadRequest)

			return
		}

		start, _ := strconv.ParseInt(q["start"][0], 10, 64)
		end, _ := strconv.ParseInt(q["end"][0], 10, 64)

		items := []WorkItem{
			{
				Issue{"XY-123", "Issue summary"},
				1643790339000, // 2022-02-02 07:25:39
				WorkItemDuration{80},
			},
		}
		rItems := []WorkItem{}
		for _, i := range items {
			if start <= i.Date && end >= i.Date {
				rItems = append(rItems, i)
			}
		}
		resp, _ := json.Marshal(rItems)

		w.Write(resp)
	})

	ytServer = &http.Server{
		Addr:    "127.0.0.1:2378",
		Handler: handler,
	}
}

func TestMain(m *testing.M) {
	now = time.Date(2022, 2, 24, 13, 13, 13, 13, time.UTC)
	home = path.Join(os.TempDir(), "ytt")
	code := m.Run()
	os.RemoveAll(home)

	os.Exit(code)
}

func Example() {
	os.Args = []string{"ytt", "init"}
	main()

	go ytServer.ListenAndServe()
	os.Args = []string{"ytt"}
	main()

	os.Args = []string{"ytt", "--start", "2022-02-03", "--end", "2022-02-22"}
	main()

	os.Args = []string{"ytt", "details"}
	main()
	ytServer.Shutdown(context.TODO())

	// Output: Created /tmp/ytt/config.json
	// 1h 20m / 143h / 159h (worked / today / month)
	// 0h / 111h (worked / month)
	// 2022-02-02	1h 20m	XY-123	Issue summary
}
