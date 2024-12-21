package main

import (
	"encoding/json"
	"github.com/showwin/speedtest-go/speedtest"
	"log"
	"os"
	"path/filepath"
	"time"
)

const (
	dataFileName = "data/metrics.json"
)

type ConnectionMetrics struct {
	Timestamp            time.Time         `json:"timestamp"`
	NetworkName          string            `json:"network_name"`
	Online               bool              `json:"online"`
	PingNanos            time.Duration     `json:"ping"`
	JitterNanos          time.Duration     `json:"jitter"`
	PacketLossPercentage float64           `json:"packet_loss"`
	DownloadMbps         float64           `json:"download"`
	UploadMbps           float64           `json:"upload"`
	Server               *speedtest.Server `json:"server"`
}

func saveMetric(metrics *ConnectionMetrics) error {
	// Create parent directories if they don't exist
	dataDirName := filepath.Dir(dataFileName)
	if err := os.MkdirAll(dataDirName, 0755); err != nil {
		log.Printf("Error creating data directory: %v\n", err)
		return err
	}

	// Open the file in append mode
	f, err := os.OpenFile(dataFileName, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Printf("Error opening data file: %v\n", err)
		return err
	}
	defer f.Close()

	jsonData, err := json.Marshal(metrics)
	if err != nil {
		log.Printf("Error encoding metrics: %v\n", err)
		return err
	}

	if _, err := f.Write(append(jsonData, '\n')); err != nil {
		log.Printf("Error writing metrics: %v\n", err)
		return err
	}

	log.Printf("Saved metrics: %s\n", jsonData)
	return nil
}

func getAllMetrics() ([]ConnectionMetrics, error) {
	f, err := os.Open(dataFileName)
	if err != nil {
		log.Printf("Error opening data file: %v\n", err)
		return nil, err
	}
	defer f.Close()

	// Read all data from the file
	var metrics []ConnectionMetrics
	decoder := json.NewDecoder(f)
	for decoder.More() {
		var m ConnectionMetrics
		if err = decoder.Decode(&m); err != nil {
			log.Printf("Error decoding metrics: %v\n", err)
			return nil, err
		}
		metrics = append(metrics, m)
	}
	return metrics, nil
}
