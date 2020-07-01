package statsagent

import (
	"time"

	"github.com/ayoubed/datadog-home-project/database"
)

// WebsiteStats contains useful metrics about website
type WebsiteStats struct {
	StatusCodeCount map[int]int
	AvgResponseTime time.Duration
	MaxResponseTime time.Duration
}

// GetStats of provided websites for a particualr span
func GetStats(urls []string, span int) map[string]WebsiteStats {
	res := database.ReadLogsForRange(urls, span)
	websitesStats := make(map[string]WebsiteStats)

	for k, v := range res {
		statusCodeCount := make(map[int]int)
		var maxResponseTime time.Duration = 0
		var sumResponseTime int64 = 0
		var avgResponseTime float64 = 0

		for _, line := range v {
			if _, ok := statusCodeCount[line.StatusCode]; ok {
				statusCodeCount[line.StatusCode]++
			} else {
				statusCodeCount[line.StatusCode] = 1
			}
			if line.LoadTime > maxResponseTime {
				maxResponseTime = line.LoadTime
			}
			sumResponseTime += int64(line.LoadTime)
		}
		if len(v) > 0 {
			avgResponseTime = float64(sumResponseTime) / float64(len(v))
		}
		websitesStats[k] = WebsiteStats{StatusCodeCount: statusCodeCount, AvgResponseTime: time.Duration(avgResponseTime), MaxResponseTime: maxResponseTime}
	}
	return websitesStats
}
