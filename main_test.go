package main

import (
	"context"
	"net/http"
	"os"
	"path"
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

		if v, ok := q["start"]; !ok || v[0] != "1643673600000" {
			http.Error(w, "Missing or invalid \"start\" query parameter", http.StatusBadRequest)

			return
		}

		if v, ok := q["end"]; !ok || v[0] != "1646092799000" {
			http.Error(w, "Missing or invalid \"end\" query parameter", http.StatusBadRequest)

			return
		}

		b, _ := os.ReadFile("response.json")
		w.Write(b)
	})

	ytServer = &http.Server{
		Addr:    "127.0.0.1:2378",
		Handler: handler,
	}
}

func TestMain(m *testing.M) {
	now = time.Date(2022, 2, 24, 13, 13, 13, 13, time.UTC)
	home = path.Join(os.TempDir(), "youtimetrack")
	code := m.Run()
	os.RemoveAll(home)

	os.Exit(code)
}

func Example() {
	os.Args = []string{"yourtimetrack", "init"}
	main()

	go ytServer.ListenAndServe()
	os.Args = []string{"yourtimetrack"}
	main()

	os.Args = []string{"youtimetrack", "details"}
	main()
	ytServer.Shutdown(context.TODO())

	// Output: Created /tmp/youtimetrack/config.json
	// 1h 20m / 143h / 159h (worked / today / month)
	// 2022-02-02 1h 20m XY-123 Issue summary
}
