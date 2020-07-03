package alerting

import (
	"fmt"
	"reflect"
	"testing"
	"time"

	"github.com/ayoubed/datadog-home-project/request"
	"github.com/ayoubed/datadog-home-project/statsagent"
)

var alertConfig AlertConfig = AlertConfig{AvailabilityInterval: 10, AvailabilityThreshold: 0.8, CheckInterval: 1}

func TestAlertLogic(t *testing.T) {
	// 0: got 0 records
	// Expected: don't send a "wesbite down" alert

	// 1: We don't have enough records on the last timeframe, availability <= threshold
	// Expected: don't send a "wesbite down" alert

	// 2: We have enough records on the last timeframe, availability > threshold, website state is up
	// Expected: don't send a "wesbite up" alert

	// 3: We have enough records on the last timeframe, availability <= threshold, website state is up
	// Expected: send a "wesbite down" alert

	// 4: We have enough records on the last timeframe, availability <= threshold, website state is down
	// Expected: dont't send a "wesbite down" alert

	// 5: We have enough records on the last timeframe, availability > threshold, website state is down
	// Expected: send a "wesbite up" alert

	start := time.Now()
	tests := []struct {
		name                 string
		URL                  string
		origin               time.Time
		records              []request.ResponseLog
		websitestateUp       bool
		expectedAlertMessage string
	}{
		{
			"0: No records",
			"https://google.com",
			time.Now(),
			[]request.ResponseLog{},
			false,
			"",
		},
		{
			"1: We don't have enough records on the last timeframe, availability <= threshold, website status up",
			"https://google.com",
			start,
			[]request.ResponseLog{
				{start.Add(-time.Duration(7) * time.Second), "200", "https://google.com", 0, 0, false},
				{start.Add(-time.Duration(6) * time.Second), "200", "https://google.com", 0, 0, true},
				{start.Add(-time.Duration(5) * time.Second), "200", "https://google.com", 0, 0, false},
				{start.Add(-time.Duration(4) * time.Second), "200", "https://google.com", 0, 0, true},
				{start.Add(-time.Duration(3) * time.Second), "200", "https://google.com", 0, 0, true},
				{start.Add(-time.Duration(2) * time.Second), "200", "https://google.com", 0, 0, true},
				{start.Add(-time.Duration(1) * time.Second), "200", "https://google.com", 0, 0, true},
			},
			true,
			"",
		},
		{
			"2: We have enough records on the last timeframe, availability > threshold, website state is up",
			"https://google.com",
			start,
			[]request.ResponseLog{
				{start.Add(-time.Duration(10) * time.Second), "200", "https://google.com", 0, 0, false},
				{start.Add(-time.Duration(8) * time.Second), "200", "https://google.com", 0, 0, true},
				{start.Add(-time.Duration(7) * time.Second), "200", "https://google.com", 0, 0, true},
				{start.Add(-time.Duration(6) * time.Second), "200", "https://google.com", 0, 0, true},
				{start.Add(-time.Duration(5) * time.Second), "200", "https://google.com", 0, 0, true},
				{start.Add(-time.Duration(4) * time.Second), "200", "https://google.com", 0, 0, true},
				{start.Add(-time.Duration(3) * time.Second), "200", "https://google.com", 0, 0, true},
				{start.Add(-time.Duration(2) * time.Second), "200", "https://google.com", 0, 0, true},
				{start.Add(-time.Duration(1) * time.Second), "200", "https://google.com", 0, 0, true},
			},
			true,
			"",
		},
		{
			"3: We have enough records on the last timeframe, availability <= threshold, website state is up",
			"https://google.com",
			start,
			[]request.ResponseLog{
				{start.Add(-time.Duration(10) * time.Second), "200", "https://google.com", 0, 0, false},
				{start.Add(-time.Duration(8) * time.Second), "200", "https://google.com", 0, 0, true},
				{start.Add(-time.Duration(7) * time.Second), "200", "https://google.com", 0, 0, false},
				{start.Add(-time.Duration(6) * time.Second), "200", "https://google.com", 0, 0, true},
				{start.Add(-time.Duration(5) * time.Second), "200", "https://google.com", 0, 0, false},
				{start.Add(-time.Duration(4) * time.Second), "200", "https://google.com", 0, 0, true},
				{start.Add(-time.Duration(3) * time.Second), "200", "https://google.com", 0, 0, true},
				{start.Add(-time.Duration(2) * time.Second), "200", "https://google.com", 0, 0, true},
				{start.Add(-time.Duration(1) * time.Second), "200", "https://google.com", 0, 0, true},
			},
			true,
			fmt.Sprintf("Website https://google.com is down. availability = 66.67%%, time = %s\n", start.Format(time.RFC1123)),
		},
		{
			"4: We have enough records on the last timeframe, availability <= threshold, website state is down",
			"https://google.com",
			start,
			[]request.ResponseLog{
				{start.Add(-time.Duration(10) * time.Second), "200", "https://google.com", 0, 0, false},
				{start.Add(-time.Duration(8) * time.Second), "200", "https://google.com", 0, 0, true},
				{start.Add(-time.Duration(7) * time.Second), "200", "https://google.com", 0, 0, false},
				{start.Add(-time.Duration(6) * time.Second), "200", "https://google.com", 0, 0, true},
				{start.Add(-time.Duration(5) * time.Second), "200", "https://google.com", 0, 0, false},
				{start.Add(-time.Duration(4) * time.Second), "200", "https://google.com", 0, 0, true},
				{start.Add(-time.Duration(3) * time.Second), "200", "https://google.com", 0, 0, true},
				{start.Add(-time.Duration(2) * time.Second), "200", "https://google.com", 0, 0, true},
				{start.Add(-time.Duration(1) * time.Second), "200", "https://google.com", 0, 0, true},
			},
			false,
			"",
		},
		{
			"5: We have enough records on the last timeframe, availability > threshold, website state is down",
			"https://google.com",
			start,
			[]request.ResponseLog{
				{start.Add(-time.Duration(10) * time.Second), "200", "https://google.com", 0, 0, false},
				{start.Add(-time.Duration(8) * time.Second), "200", "https://google.com", 0, 0, true},
				{start.Add(-time.Duration(7) * time.Second), "200", "https://google.com", 0, 0, true},
				{start.Add(-time.Duration(6) * time.Second), "200", "https://google.com", 0, 0, true},
				{start.Add(-time.Duration(5) * time.Second), "200", "https://google.com", 0, 0, true},
				{start.Add(-time.Duration(4) * time.Second), "200", "https://google.com", 0, 0, true},
				{start.Add(-time.Duration(3) * time.Second), "200", "https://google.com", 0, 0, true},
				{start.Add(-time.Duration(2) * time.Second), "200", "https://google.com", 0, 0, true},
				{start.Add(-time.Duration(1) * time.Second), "200", "https://google.com", 0, 0, true},
			},
			false,
			fmt.Sprintf("Website https://google.com is up. availability = 88.89%%, time = %s\n", start.Format(time.RFC1123)),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			availability := statsagent.GetAvailabilityForRecords(tt.records, tt.origin)
			alertMessage := getAlertMessage(tt.origin, tt.URL, tt.websitestateUp, 1, availability, alertConfig)

			if !reflect.DeepEqual(alertMessage, tt.expectedAlertMessage) {
				t.Errorf("Got %v, want %v", alertMessage, tt.expectedAlertMessage)
			}
		})
	}
}
