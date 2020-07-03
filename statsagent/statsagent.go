package statsagent

import (
	"time"

	"github.com/ayoubed/datadog-home-project/database"
	"github.com/ayoubed/datadog-home-project/request"
)

// WebsiteStats contains useful metrics about website
type WebsiteStats struct {
	StatusCodeCount    map[int]int
	AvgResponseTime    time.Duration
	MaxResponseTime    time.Duration
	AvgTimeToFirstByte time.Duration
	MaxTimeToFirstByte time.Duration
	Availability       float64
}

// GetStats of provided websites for a particualr span
func GetStats(urls []string, span int) map[string]WebsiteStats {
	res := database.ReadLogsForRange(urls, span)
	websitesStats := make(map[string]WebsiteStats)

	for k, v := range res {
		statusCodeCount := make(map[int]int)
		var sumResponseTime int64 = 0
		var maxResponseTime time.Duration = 0
		var avgResponseTime float64 = 0
		var sumTimeToFirstByte int64 = 0
		var maxTimeToFirstByte time.Duration = 0
		var avgTimeToFirstByte float64 = 0
		var successCount float64 = 0
		var availability float64 = 0

		for _, line := range v {
			if line.Success {
				successCount++
			}
			if _, ok := statusCodeCount[line.StatusCode]; ok {
				statusCodeCount[line.StatusCode]++
			} else {
				statusCodeCount[line.StatusCode] = 1
			}
			if line.LoadTime > maxResponseTime {
				maxResponseTime = line.LoadTime
			}
			if line.TTFB > maxTimeToFirstByte {
				maxTimeToFirstByte = line.TTFB
			}
			sumResponseTime += int64(line.LoadTime)
			sumTimeToFirstByte += int64(line.TTFB)
		}
		if len(v) > 0 {
			avgResponseTime = float64(sumResponseTime) / float64(len(v))
			avgTimeToFirstByte = float64(sumTimeToFirstByte) / float64(len(v))
			availability = successCount / float64(len(v))
		}
		websitesStats[k] = WebsiteStats{StatusCodeCount: statusCodeCount, AvgResponseTime: time.Duration(avgResponseTime), MaxResponseTime: maxResponseTime, AvgTimeToFirstByte: time.Duration(avgTimeToFirstByte), MaxTimeToFirstByte: maxTimeToFirstByte, Availability: availability}
	}
	return websitesStats
}

type AvailabilityRange struct {
	Availability float64
	Start        time.Time
	Records      []request.ResponseLog
}

// GetRangeAvailability : get availability for a particular timeframe
func GetRangeAvailability(urls []string, t time.Time, timeframe int) map[string]AvailabilityRange {
	var start time.Time = time.Now()
	res := database.ReadLogsForRange2(urls, t, timeframe)
	websitesAvailability := make(map[string]AvailabilityRange)

	for k, v := range res {
		var successCount float64 = 0
		var availability float64 = 0

		for _, line := range v {
			if line.Success {
				successCount++
			}
		}

		if len(v) > 0 {
			availability = successCount / float64(len(v))
			start = v[0].Timestamp
		}
		websitesAvailability[k] = AvailabilityRange{Availability: availability, Start: start, Records: v}
	}
	return websitesAvailability
}
