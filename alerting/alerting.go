package alerting

import (
	"time"

	"github.com/ayoubed/datadog-home-project/statsagent"
	"github.com/fatih/color"
)

var red *color.Color = color.New(color.FgRed)
var green *color.Color = color.New(color.FgGreen)

// websiteUp is a map that keeps track of the state of each website we're monitoring
var websiteUp map[string]bool = make(map[string]bool)

// AlertConfig represents all the useful info for our alert logic
type AlertConfig struct {
	AvailabilityInterval  int64   `json:"availabilityInterval"`
	AvailabilityThreshold float64 `json:"availabilityThreshold"`
	CheckInterval         int     `json:"checkInterval"`
}

// Run monitors the availability of websites
// It send an alert to the dashboard, if the availability of some website over a given interval
// is bellow the given the threshold
func Run(alertc chan string, websitesMap map[string]int64, alertConfig AlertConfig) {
	urls := make([]string, 0)
	for k := range websitesMap {
		websiteUp[k] = true
		urls = append(urls, k)
	}
	ticker := time.NewTicker(time.Duration(alertConfig.CheckInterval) * time.Second)

	for {
		select {
		case t := <-ticker.C:
			for _, url := range urls {
				v := statsagent.GetAvailabilityForTimeFrame(url, t, alertConfig.AvailabilityInterval)
				var tm int64 = (v.Start.Unix() - (t.Unix() - alertConfig.AvailabilityInterval))

				if tm >= 0 && tm <= websitesMap[url] && (v.Availability <= alertConfig.AvailabilityThreshold && websiteUp[url] == true) || (v.Availability > alertConfig.AvailabilityThreshold && websiteUp[url] == false) {
					var alertMessage string

					if v.Availability > alertConfig.AvailabilityThreshold {
						alertMessage = green.Sprintf("Website %v is up. availability = %.2f%%, time = %s\n", url, 100*v.Availability, t.Format(time.RFC1123))
					} else {
						alertMessage = red.Sprintf("Website %v is down. availability = %.2f%%, time = %s\n", url, 100*v.Availability, t.Format(time.RFC1123))
					}

					alertc <- alertMessage
					websiteUp[url] = v.Availability > alertConfig.AvailabilityThreshold
				}
			}
		}
	}
}
