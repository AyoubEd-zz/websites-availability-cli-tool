package request

import (
	"net"
	"net/http"
	"net/http/httptrace"
	"strconv"

	"time"
)

// ResponseLog represents the info we keep as log from the requests we issue
type ResponseLog struct {
	Timestamp  time.Time
	StatusCode string
	URL        string
	TTFB       time.Duration
	LoadTime   time.Duration
	Success    bool
}

// Send performs a request to the given URL
func Send(t time.Time, url string) (ResponseLog, error) {
	var (
		start time.Time
		ttfb  time.Duration
	)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return ResponseLog{}, err
	}
	client := &http.Client{
		Timeout: 15 * time.Second,
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
		if err, ok := err.(net.Error); ok && err.Timeout() {
			return ResponseLog{t, err.Error(), url, time.Duration(0), time.Duration(0), false}, nil
		}
		return ResponseLog{}, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return ResponseLog{t, strconv.Itoa(resp.StatusCode), url, time.Duration(0), time.Duration(0), false}, nil
	}
	return ResponseLog{t, strconv.Itoa(resp.StatusCode), url, ttfb, time.Since(start), true}, nil
}
