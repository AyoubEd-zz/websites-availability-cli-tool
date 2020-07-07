package database

import (
	"fmt"
	"time"

	"github.com/ayoubed/datadog-home-project/request"
	"github.com/influxdata/influxdb/client/v2"
)

var (
	dbName InfluxDb
)

// Database interface abstracts the interactions woth influxDB
// and provides and database agnostic interface for use
type Database interface {
	Initialize() error
	GetDatabaseName() string
	AddResponseLog(responseLog request.ResponseLog) error
	GetRangeRecords(span int) []client.Result
}

// Type is the database type
// ps: if we decide we want to use another db we can add it here
type Type struct {
	InfluxDb InfluxDb `json:"influxDb"`
}

//Set sets the database name, and initializes the database
func Set(database Type) error {
	dbName = database.InfluxDb
	if err := dbName.Initialize(); err != nil {
		return err
	}
	return nil
}

// WriteLogToDB writes logs to our database
func WriteLogToDB(responseLog request.ResponseLog) error {
	if err := dbName.AddRecord(responseLog); err != nil {
		return fmt.Errorf("error while writing a log to the database:\n %v", err)
	}
	return nil
}

// GetRecordsForURL gets records from the database for all the given URLs
// the records timestamp is bounded between and [origin - timeframe, origin]
func GetRecordsForURL(url string, origin time.Time, timeframe int64) ([]request.ResponseLog, error) {
	res, err := dbName.GetRecordsForURL(url, origin, timeframe)
	if err != nil {
		return nil, err
	}
	return res, nil
}
