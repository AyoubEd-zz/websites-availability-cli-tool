package request

import (
	"fmt"
	"net/http"
	"net/http/httptrace"

	"time"
)

// ResponseLog represents the info we keep as log from the requests we issue
type ResponseLog struct {
	Timestamp  time.Time
	StatusCode int
	URL        string
	TTFB       time.Duration
	LoadTime   time.Duration
	Success    bool
}

// Send performs a request to the given URL
func Send(t time.Time, url string, logc chan ResponseLog) error {
	var (
		start time.Time
		ttfb  time.Duration
	)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return err
	}
	client := &http.Client{
		Timeout: 10 * time.Second,
	}

	trace := &httptrace.ClientTrace{
		GotFirstResponseByte: func() {
			ttfb = time.Since(start)
		},
	}
	req = req.WithContext(httptrace.WithClientTrace(req.Context(), trace))
	start = time.Now()

	resp, err := client.Do(req)
	if err != nil {
		fmt.Printf("%+v", err)
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		logc <- ResponseLog{t, resp.StatusCode, url, ttfb, time.Since(start), false}
	} else {
		logc <- ResponseLog{t, resp.StatusCode, url, ttfb, time.Since(start), true}
	}

	return nil
}
