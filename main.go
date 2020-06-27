package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"time"

	"github.com/ayoubed/datadog-home-project/request"
)

// Website representes the entities we want to monitor
type Website struct {
	URL           string `json:"url"`
	CheckInterval int    `json:"checkInterval"`
}

func main() {
	websites, err := getWebsites()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error loading website config: %v\n", err)
		os.Exit(1)
	}

	if err := runMonitor(websites); err != nil {
		fmt.Fprintf(os.Stderr, "The website monitor encountered an error: %v\n", err)
		os.Exit(1)
	}
}

func getWebsites() ([]Website, error) {
	configFile, err := os.Open("config.json")
	if err != nil {
		return nil, fmt.Errorf("%v", err)
	}
	defer configFile.Close()

	var websites []Website
	byteContent, err := ioutil.ReadAll(configFile)
	if err != nil {
		return nil, fmt.Errorf("%v", err)
	}

	json.Unmarshal(byteContent, &websites)

	return websites, nil
}

func runMonitor(websites []Website) error {
	// Start goroutines to ping websites
	done := make(chan bool, 1)
	errc := make(chan error)
	defer close(done)

	for _, ws := range websites {
		go startTicker(ws, done, errc)
	}

	for {
		select {
		case err := <-errc:
			return err
		case <-done:
			close(errc)
		}
	}

}

func startTicker(website Website, done chan bool, errc chan error) {
	ticker := time.NewTicker(time.Duration(website.CheckInterval) * time.Millisecond)
	for {
		select {
		case <-done:
			ticker.Stop()
			return
		case t := <-ticker.C:
			if err := request.Send(t, website.URL); err != nil {
				errc <- err
			}
		}
	}
}
