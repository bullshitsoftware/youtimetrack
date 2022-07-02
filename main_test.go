package main

import (
	"context"
	"net/http"
	"os"
	"path"
	"testing"
	"time"
)

func TestMain(m *testing.M) {
	now = time.Date(2022, 2, 24, 13, 13, 13, 13, time.Now().Location())
	home = path.Join(os.TempDir(), "youtimetrack")
	code := m.Run()
	os.RemoveAll(home)

	os.Exit(code)
}

func Example() {
	os.Args = []string{"yourtimetrack", "init"}
	main()

	http.HandleFunc("/api/workItems", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("[{\"duration\":{\"minutes\":80}}]"))
	})
	server := &http.Server{Addr: "127.0.0.1:2378"}
	go server.ListenAndServe()

	os.Args = []string{"yourtimetrack"}
	main()
	server.Shutdown(context.TODO())

	// Output: Created /tmp/youtimetrack/config.json
	// 1h 20m / 143h / 159h (worked / today / month)
}
