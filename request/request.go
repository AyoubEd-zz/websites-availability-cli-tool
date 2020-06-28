package request

import (
	"fmt"
	"net/http"
	"net/http/httptrace"

	"time"
)

// ResponseLog represents the info we keep as log from the requests we issue
type ResponseLog struct {
	Timestamp  int64
	StatusCode int
	URL        string
	TTFB       time.Duration
	LoadTime   time.Duration
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

	trace := &httptrace.ClientTrace{
		GotFirstResponseByte: func() {
			ttfb = time.Since(start)
		},
	}
	req = req.WithContext(httptrace.WithClientTrace(req.Context(), trace))
	start = time.Now()

	resp, err := http.DefaultTransport.RoundTrip(req)
	if err != nil {
		fmt.Println(err)
	} else {
		logc <- ResponseLog{start.Unix(), resp.StatusCode, url, ttfb, time.Since(start)}
	}

	return nil
}
