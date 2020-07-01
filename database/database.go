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

// ReadLogsPeriodically reads logs from the database
func ReadLogsPeriodically(urls []string, interval, span int) {
	ticker := time.NewTicker(time.Duration(interval) * time.Second)

	for {
		select {
		case t := <-ticker.C:
			fmt.Printf("[%vs stats] ------------------%v------------------\n", interval, t)
			for _, url := range urls {
				fmt.Printf("> %v\n", url)
				res := dbName.GetRangeRecords(url, span)
				for _, line := range res {
					fmt.Printf("%+v\n", line)
				}
			}
		}
	}
}
