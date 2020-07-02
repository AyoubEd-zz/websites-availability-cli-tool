package database

import (
	"fmt"
	"os"

	"github.com/ayoubed/datadog-home-project/request"
	"github.com/influxdata/influxdb/client/v2"
)

var (
	dbName InfluxDb
)

// Database interface
type Database interface {
	Initialize() error
	GetDatabaseName() string
	AddResponseLog(responseLog request.ResponseLog) error
	GetRangeRecords(span int) []client.Result
}

// Type is the database type
type Type struct {
	InfluxDb InfluxDb `json:"influxDb"`
}

//Set sets the database name
func Set(database Type) {
	dbName = database.InfluxDb
	if err := dbName.Initialize(); err != nil {
		fmt.Println("Failed to Intialize Database ")
		os.Exit(3)
	}
}

// WriteLogToDB writes logs to our database
func WriteLogToDB(responseLog request.ResponseLog) {
	go dbName.AddResponseLog(responseLog)
}

// ReadLogsForRange reads logs from the database
func ReadLogsForRange(urls []string, span int) map[string][]request.ResponseLog {
	logsForURL := make(map[string][]request.ResponseLog)
	for _, url := range urls {
		logsForURL[url] = dbName.GetRangeRecords(url, span)
	}
	return logsForURL
}
