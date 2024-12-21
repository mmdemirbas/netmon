package main

import (
	"github.com/showwin/speedtest-go/speedtest"
	"log"
	"time"
)

const (
	interval = 1 * time.Minute // Interval between measurements
)

func startCollector() {
	log.Printf("Starting data collection every %v\n", interval)
	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	collect()
	for range ticker.C {
		collect()
	}
}

// collect performs a speed test, gets the network name, and saves the metrics.
func collect() {
	log.Printf("Starting speed test...\n")

	serverList, _ := speedtest.FetchServers()
	targets, _ := serverList.FindServer([]int{})

	// Get the network name
	networkName, err := getNetworkName()
	if err != nil {
		log.Printf("Error getting network name: %v", err)
		networkName = "Unknown" // Set a default value on error
	}

	log.Printf("Starting speed test on network: %s\n", networkName)

	// Assume offline initially
	online := false
	defer func() {
		// If no server was tested, or any error occurred, the connection is considered offline
		if !online {
			metrics := &ConnectionMetrics{
				Timestamp:   time.Now(),
				NetworkName: networkName,
				Online:      false,
			}
			if err = saveMetric(metrics); err != nil {
				log.Printf("Error saving offline data: %v\n", err)
			} else {
				log.Printf("Saved offline data: %v\n", metrics)
			}
		}
	}()

	for _, server := range targets {
		log.Printf("Testing server: %s\n", server.Name)
		err = server.TestAll()
		if err != nil {
			log.Printf("Error testing server: %v\n", err)
			continue // Try the next server
		}

		online = true
		metrics := &ConnectionMetrics{
			Timestamp:            time.Now(),
			NetworkName:          networkName,
			Online:               true,
			PingNanos:            server.Latency,
			JitterNanos:          server.Jitter,
			PacketLossPercentage: server.PacketLoss.LossPercent(),
			DownloadMbps:         server.DLSpeed.Mbps(),
			UploadMbps:           server.ULSpeed.Mbps(),
			Server:               server,
		}

		err = saveMetric(metrics)
		if err != nil {
			log.Printf("Error saving online data: %v\n", err)
		} else {
			log.Printf("Saved online data: %v\n", metrics)
		}
		return
	}

	log.Printf("Error on speed test: %v\n", err)
}
