package alerting

import (
	"time"

	"github.com/ayoubed/datadog-home-project/dashboard"
	"github.com/ayoubed/datadog-home-project/statsagent"
)

var websiteUp map[string]bool = make(map[string]bool)

// Run monitors the availability of the websites
func Run(websitesMap map[string]int64, threshold float64) {
	// every 5 seconds check the availablity of websites
	urls := make([]string, 0)
	for k := range websitesMap {
		websiteUp[k] = true
		urls = append(urls, k)
	}
	ticker := time.NewTicker(5 * time.Second)

	for {
		select {
		case t := <-ticker.C:
			res := statsagent.GetRangeAvailability(urls, t, 120)
			for k, v := range res {
				var tm int64 = (v.Start.Unix() - (t.Unix() - 120)) * 1000
				if tm >= 0 && tm <= websitesMap[k] && (v.Availability <= threshold && websiteUp[k] == true) || (v.Availability > threshold && websiteUp[k] == false) {
					dashboard.Alert(t, k, v.Availability, v.Availability > threshold, v.Start)
					// fmt.Println("aleeeeeeeeeeeeert", k, t.Unix(), v.Start.Unix(), tm, websitesMap[k])
					websiteUp[k] = v.Availability > threshold
				} else {
					// fmt.Println("> ", k, t.Format("15:04:05"), v.Start.Format("15:04:05"), tm, websitesMap[k])
				}
			}
		}
	}
}
