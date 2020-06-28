package statsagent

import (
	"fmt"
	"sync"
	"time"

	"github.com/ayoubed/datadog-home-project/request"
)

type WebsiteStat struct {
	availability       int
	responseCodeCount  map[int]int
	avgResponseTime    time.Duration
	maxResponseTime    time.Duration
	avgTimeToFirstByte time.Duration
	maxTimeToFirstByte time.Duration
}

type websiteStatsQueue struct {
	refreshInterval int64
	queue           []request.ResponseLog
	mux             sync.Mutex
}

func ProcessLogs(logc chan request.ResponseLog, errc chan error) {
	done := make(chan bool, 1)

	ws1, ws2 := websiteStatsQueue{refreshInterval: 10, queue: []request.ResponseLog{}}, websiteStatsQueue{refreshInterval: 60, queue: []request.ResponseLog{}}
	go ws1.removeOutdatedLogs(done, errc)
	go ws2.removeOutdatedLogs(done, errc)
	go ws1.startTicker(done, errc)
	go ws2.startTicker(done, errc)
	for log := range logc {
		ws1.Add(log)
		ws2.Add(log)
	}
}

func (c *websiteStatsQueue) Add(log request.ResponseLog) {
	c.mux.Lock()
	defer c.mux.Unlock()
	c.queue = append(c.queue, log)
}

func (c *websiteStatsQueue) removeOutdatedLogs(done chan bool, errc chan error) {
	timer := time.NewTicker(time.Duration(10) * time.Second)

	for {
		select {
		case t := <-timer.C:
			c.mux.Lock()
			newQueue := make([]request.ResponseLog, 0)
			for _, lg := range c.queue {
				if t.Unix()-lg.Timestamp <= 10 {
					newQueue = append(newQueue, lg)
				}
			}
			c.queue = newQueue
			c.mux.Unlock()
		case <-done:
			timer.Stop()
			return
		}
	}
}

func (c *websiteStatsQueue) startTicker(done chan bool, errc chan error) {
	timer := time.NewTicker(time.Duration(c.refreshInterval) * time.Second)

	for {
		select {
		case t := <-timer.C:
			c.mux.Lock()
			fmt.Println(t, "--------------------------------------------", c.refreshInterval)
			for _, lg := range c.queue {
				fmt.Println(time.Unix(lg.Timestamp, 0), lg.URL)
			}
			c.mux.Unlock()
		case <-done:
			timer.Stop()
		}
	}
}
