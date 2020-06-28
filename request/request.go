package request

import (
	"net/http"
	"net/http/httptrace"

	"time"
)

// ResponseLog represents the info we keep as log from the requests we issue
type ResponseLog struct {
	Timestamp  int64
	statusCode int
	URL        string
	ttfb       time.Duration
	loadTime   time.Duration
}

// Send performs a request to the given URL
func Send(t time.Time, url string, logc chan ResponseLog) error {
	var (
		start  time.Time
		ttfb   time.Duration
		resLog ResponseLog
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
		return err
	}

	resLog = ResponseLog{start.Unix(), resp.StatusCode, url, ttfb, time.Since(start)}
	logc <- resLog
	return nil
}
