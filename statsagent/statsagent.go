package statsagent

import (
	"fmt"
	"sort"
	"sync"
	"time"

	"github.com/ayoubed/datadog-home-project/request"
)

type WebsiteStat struct {
	numberOfLogs       int64
	availability       int
	responseCodeCount  map[int]int
	avgResponseTime    int64
	maxResponseTime    time.Duration
	sumResponseTime    time.Duration
	avgTimeToFirstByte int64
	maxTimeToFirstByte time.Duration
	sumTimeToFirstByte time.Duration
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

	mp := map[string]*WebsiteStat{}
	for _, lg := range sortedLogs {
		if v, ok := mp[lg.URL]; ok {
			v.numberOfLogs++
			if lg.TTFB > v.maxTimeToFirstByte {
				v.maxTimeToFirstByte = lg.TTFB
			}
			if lg.LoadTime > v.maxResponseTime {
				v.maxResponseTime = lg.LoadTime
			}
			v.responseCodeCount[lg.StatusCode]++

			v.sumResponseTime += lg.LoadTime
			v.avgResponseTime = int64(v.sumResponseTime) / v.numberOfLogs

			v.sumTimeToFirstByte += lg.LoadTime
			v.avgTimeToFirstByte = int64(v.sumResponseTime) / v.numberOfLogs
		} else {
			mp[lg.URL] = &WebsiteStat{1, 0, map[int]int{lg.StatusCode: 1}, lg.LoadTime.Milliseconds(), lg.LoadTime, lg.LoadTime, lg.TTFB.Milliseconds(), lg.TTFB, lg.TTFB}
		}
	}

	for k, v := range mp {
		fmt.Println(k, v)
	}
}
