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

	serverList, _ := speedtest.FetchServers()
	targets, _ := serverList.FindServer([]int{})

	// Assume offline initially
	online := false
	defer func() {
		// If no server was tested, or any error occurred, the connection is considered offline
		if !online {
			metrics := &Metrics{
				Timestamp:   time.Now(),
				NetworkName: networkName,
				Online:      false,
			}
			if err = saveMetric(metrics); err != nil {
				logger.Errorf("Error saving offline data: %v", err)
			}
		}
	}()

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
