package database

import (
	"fmt"
	"os"
	"time"

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

// GetRecordsForURL gets records from the database for all the given URLs
// the records timestamp is bounded between and [origin - timeframe, origin]
func GetRecordsForURL(url string, origin time.Time, timeframe int64) []request.ResponseLog {
	return dbName.GetRecordsForURL(url, origin, timeframe)
}
