package main

import (
	"github.com/showwin/speedtest-go/speedtest"
	"time"
)

func startCollector(collectorInterval *time.Duration) *time.Ticker {
	ticker := time.NewTicker(*collectorInterval)
	go func() {
		collect()
		for range ticker.C {
			collect()
		}
	}()
	return ticker
}

// collect performs a speed test, gets the network name, and saves the metrics.
func collect() {
	// Get the network name
	networkName, err := getNetworkName()
	if err != nil {
		logger.Errorf("Error getting network name: %v", err)
		networkName = "Unknown" // Set a default value on error
	}

	logger.Infof("Starting speed test on network: %s", networkName)

	// Assume offline initially; set online=true only after a successful server test.
	online := false
	defer func() {
		if !online {
			metrics := &Metrics{
				Timestamp:   time.Now(),
				NetworkName: networkName,
				Online:      false,
			}
			if err = saveMetric(metrics); err != nil {
				logger.Errorf("error saving offline metrics: %v", err)
			}
		}
	}()

	serverList, err := speedtest.FetchServers()
	if err != nil {
		logger.Errorf("failed to fetch speedtest servers: %v", err)
		return
	}
	targets, err := serverList.FindServer([]int{})
	if err != nil {
		logger.Errorf("failed to find speedtest servers: %v", err)
		return
	}

	for _, server := range targets {
		err = server.TestAll()
		if err != nil {
			logger.Warningf("Error testing server %s: %v", server.Name, err)
			continue // Try the next server
		}

		online = true
		metrics := &Metrics{
			Timestamp:   time.Now(),
			NetworkName: networkName,
			Online:      true,
			Data:        server,
		}

		err = saveMetric(metrics)
		if err != nil {
			logger.Errorf("Error saving online data: %v", err)
		}
		return
	}

	logger.Errorf("Error on speed test: %v", err)
}
