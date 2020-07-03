package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"

	"github.com/ayoubed/datadog-home-project/alerting"
	"github.com/ayoubed/datadog-home-project/dashboard"
	"github.com/ayoubed/datadog-home-project/database"
	"github.com/ayoubed/datadog-home-project/monitor"
	"github.com/ayoubed/datadog-home-project/request"
	"golang.org/x/sync/errgroup"
)

// Config struct containing websites config(url, check interval), database data(host, dbaname, username, password)
type Config struct {
	Websites  []monitor.Website    `json:"websites"`
	Database  database.Type        `json:"database"`
	Dashboard []dashboard.View     `json:"dashboard"`
	Alert     alerting.AlertConfig `json:"alerting"`
}

func main() {
	config, err := getConfig()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error loading config: %v\n", err)
		os.Exit(1)
	}

	if err := database.Set(config.Database); err != nil {
		fmt.Fprintf(os.Stderr, "Error setting up the database: %v\n", err)
		os.Exit(1)
	}

	websiteList := []string{}
	websiteMap := make(map[string]int64)
	for _, ws := range config.Websites {
		websiteList = append(websiteList, ws.URL)
		websiteMap[ws.URL] = int64(ws.CheckInterval)
	}

	ctx, done := context.WithCancel(context.Background())
	g, gctx := errgroup.WithContext(ctx)

	// Start goroutines to ping websites
	logc := make(chan request.ResponseLog)
	alertc := make(chan string)
	defer close(logc)
	defer close(alertc)

	go dashboard.Run(websiteList, config.Dashboard, alertc, done)
	go alerting.Run(alertc, websiteMap, config.Alert)

	g.Go(func() error {
		return monitor.ProcessLogs(gctx, logc)
	})

	for _, ws := range config.Websites {
		ws := ws
		g.Go(func() error {
			return monitor.StartWebsiteMonitor(gctx, ws, logc)
		})
	}

	if err := g.Wait(); err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		os.Exit(1)
	}

	// if err := monitor.Run(config.Websites, logc, errc, done); err != nil {
	// 	fmt.Fprintf(os.Stderr, "The website monitor encountered an error: %v\n", err)
	// 	os.Exit(1)
	// }
}

func getConfig() (Config, error) {
	configFile, err := os.Open("config.json")
	if err != nil {
		return Config{}, fmt.Errorf("%v", err)
	}
	defer configFile.Close()

	configByteContent, err := ioutil.ReadAll(configFile)
	if err != nil {
		return Config{}, err
	}

	var config Config

	if err := json.Unmarshal(configByteContent, &config); err != nil {
		return Config{}, err
	}

	return config, nil
}
