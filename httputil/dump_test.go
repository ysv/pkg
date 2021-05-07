package httputil

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestDumpResponse(t *testing.T) {
	//const body = "Go is a general-purpose language designed with systems programming in mind."
	const body = `{"a": "b"}`
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Date", "Wed, 19 Jul 1972 19:00:00 GMT")
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprintln(w, body)
	}))
	defer ts.Close()

	//resp, err := http.Post(fmt.Sprintf("%v/%v", ts.URL, "foo"), "application/json", strings.NewReader(`{"ab": "cd"}`))
	ctx, cancelF := context.WithCancel(context.Background())

	req, err := http.NewRequestWithContext(ctx, "GET", "https://api.test.sologenic.org/api/v1/ohlc?symbol=534F4C4F00000000000000000000000000000000%2BrsoLo2S1kiGeCcn6hCUXVrCpGMWLrRrLZz%2FXRP&from=1615220592&to=1620404592&period=1h", nil)
	go func() {
		<-time.After(time.Millisecond * 10)
		cancelF()
	}()
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()

	dump, err := DumpResponse(resp)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(dump)

	fmt.Println("dump2 ------------------------------------------------")
	dump, err = DumpResponse(resp)
	if err != nil {
		fmt.Println("Error: not nil")
		log.Fatal(err)
	}

	fmt.Println("Error: nil")
	fmt.Println(string(dump))
}
