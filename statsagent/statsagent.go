package statsagent

import (
	"github.com/ayoubed/datadog-home-project/database"
	"github.com/ayoubed/datadog-home-project/request"
)

// ProcessLogs receives logs from their dedicated channel, processes them and sends to them to be written to the database
func ProcessLogs(logc chan request.ResponseLog, errc chan error) {
	for log := range logc {
		database.WriteLogToDB(log)
	}
}
