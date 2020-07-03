package monitor

import (
	"context"
	"fmt"
	"time"

	"github.com/ayoubed/datadog-home-project/database"
	"github.com/ayoubed/datadog-home-project/request"
)

// Website representes the entities we want to monitor
type Website struct {
	URL           string `json:"url"`
	CheckInterval int    `json:"checkInterval"`
}

// StartWebsiteMonitor starts a ticker for the given website
// it sends a request following a user-defined interval
func StartWebsiteMonitor(ctx context.Context, website Website, logc chan request.ResponseLog) error {
	ticker := time.NewTicker(time.Duration(website.CheckInterval) * time.Second)
	for {
		select {
		case <-ctx.Done():
			ticker.Stop()
			return nil
		case t := <-ticker.C:
			log, err := request.Send(t, website.URL)
			if err != nil {
				return fmt.Errorf("error while monitoring %v:\n Details: %v", website, err)
			}
			logc <- log
		}
	}
}

// ProcessLogs reads logs from the log channel and processes them
// in our case we write logs in our database
func ProcessLogs(ctx context.Context, logc chan request.ResponseLog) error {
	for {
		select {
		case <-ctx.Done():
			return nil
		case log := <-logc:
			if err := database.WriteLogToDB(log); err != nil {
				return fmt.Errorf("error while processing a log:\n %v", err)
			}
		}
	}
}
