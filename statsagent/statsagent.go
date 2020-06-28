package statsagent

import (
	"fmt"
	"sort"
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
	ID              int
	refreshInterval int64
	displayInterval int64
	queue           []request.ResponseLog
	mux             sync.Mutex
}

// ProcessLogs receives logs from the ping agent andf stores the latest logs in memory
func ProcessLogs(logc chan request.ResponseLog, errc chan error) {
	done := make(chan bool, 1)

	ws1, ws2 := websiteStatsQueue{ID: 1, refreshInterval: 60, displayInterval: 10, queue: []request.ResponseLog{}}, websiteStatsQueue{ID: 2, refreshInterval: 3600, displayInterval: 60, queue: []request.ResponseLog{}}
	go ws1.removeOutdatedLogs(done, errc)
	go ws2.removeOutdatedLogs(done, errc)
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
	var ticks int64 = 0

	for {
		select {
		case t := <-timer.C:
			ticks++
			c.mux.Lock()
			newQueue := make([]request.ResponseLog, 0)
			for _, lg := range c.queue {
				if t.Unix()-lg.Timestamp <= c.refreshInterval {
					newQueue = append(newQueue, lg)
				}
			}
			c.queue = newQueue
			if (ticks*10)%c.displayInterval == 0 {
				go sendUpdatedStats(c.ID, c.queue, t, c.refreshInterval/6)
			}
			c.mux.Unlock()
		case <-done:
			timer.Stop()
			return
		}
	}
}

func sendUpdatedStats(ID int, q []request.ResponseLog, t time.Time, inter int64) {
	sortedLogs := q
	sort.Slice(sortedLogs, func(a, b int) bool {
		return sortedLogs[a].Timestamp < sortedLogs[b].Timestamp
	})
	fmt.Println("[", ID, "]", t, "--------------------------------------------", inter)
	for _, lg := range sortedLogs {
		fmt.Println(time.Unix(lg.Timestamp, 0), lg.URL)
	}
}
