package statsagent

import (
	"fmt"
	"time"

	"github.com/ayoubed/datadog-home-project/database"
	"github.com/ayoubed/datadog-home-project/request"
)

// WebsiteStats contains useful metrics about website
type WebsiteStats struct {
	StatusCodeCount    map[string]int
	AvgResponseTime    time.Duration
	MaxResponseTime    time.Duration
	AvgTimeToFirstByte time.Duration
	MaxTimeToFirstByte time.Duration
	Availability       float64
}

// GetStats of provided websites for a particular timeframe
func GetStats(urls []string, origin time.Time, timeframe int64) map[string]WebsiteStats {
	websitesStats := make(map[string]WebsiteStats)

	for _, url := range urls {
		records, err := database.GetRecordsForURL(url, origin, timeframe)
		if err != nil {
			// todo
		}
		statusCodeCount := make(map[string]int)
		var sumResponseTime int64 = 0
		var maxResponseTime time.Duration = 0
		var avgResponseTime float64 = 0
		var sumTimeToFirstByte int64 = 0
		var maxTimeToFirstByte time.Duration = 0
		var avgTimeToFirstByte float64 = 0
		var successCount float64 = 0
		var availability float64 = 0

		for _, line := range records {
			if _, ok := statusCodeCount[line.StatusCode]; ok {
				statusCodeCount[line.StatusCode]++
			} else {
				statusCodeCount[line.StatusCode] = 1
			}

			if line.Success {
				successCount++
				if line.LoadTime > maxResponseTime {
					maxResponseTime = line.LoadTime
				}
				if line.TTFB > maxTimeToFirstByte {
					maxTimeToFirstByte = line.TTFB
				}
				sumResponseTime += int64(line.LoadTime)
				sumTimeToFirstByte += int64(line.TTFB)
			}
		}
		if successCount > 0 {
			avgResponseTime = float64(sumResponseTime) / float64(successCount)
			avgTimeToFirstByte = float64(sumTimeToFirstByte) / float64(successCount)
			availability = successCount / float64(len(records))
		}
		websitesStats[url] = WebsiteStats{StatusCodeCount: statusCodeCount, AvgResponseTime: time.Duration(avgResponseTime), MaxResponseTime: maxResponseTime, AvgTimeToFirstByte: time.Duration(avgTimeToFirstByte), MaxTimeToFirstByte: maxTimeToFirstByte, Availability: availability}
	}
	return websitesStats
}

type AvailabilityRange struct {
	Availability float64
	Start        time.Time
}

// GetAvailabilityForTimeFrame computes the availability of a Website
// given a time origin and a timeframe
func GetAvailabilityForTimeFrame(url string, origin time.Time, timeframe int64) (AvailabilityRange, error) {
	records, err := database.GetRecordsForURL(url, origin, timeframe)
	if err != nil {
		return AvailabilityRange{}, fmt.Errorf("error while calculating availabililty for %v: %v", url, err)
	}
	return GetAvailabilityForRecords(records, origin), nil
}

// GetAvailabilityForRecords returns the availability given a slice of records
func GetAvailabilityForRecords(records []request.ResponseLog, origin time.Time) AvailabilityRange {
	var start time.Time = origin
	var successCount float64 = 0
	var availability float64 = 0

	for _, line := range records {
		if line.Success {
			successCount++
		}
	}

	if successCount > 0 {
		availability = successCount / float64(len(records))
		start = records[0].Timestamp
	}
	return AvailabilityRange{Availability: availability, Start: start}
}
