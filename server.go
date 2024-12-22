package main

import (
	"embed"
	"encoding/json"
	"fmt"
	"io/fs"
	"net/http"
)

// Embed the 'static' directory which contains the HTML and JS files
//
//go:embed static
var staticFiles embed.FS

type MetricsDto struct {
	EpochMillis          int64   `json:"EpochMillis"`
	NetworkName          string  `json:"NetworkName"`
	IsOnline             bool    `json:"IsOnline"`
	PingMillis           int64   `json:"PingMillis"`
	JitterMillis         int64   `json:"JitterMillis"`
	PacketLossPercentage float64 `json:"PacketLossPercentage"`
	DownloadMbps         float64 `json:"DownloadMbps"`
	UploadMbps           float64 `json:"UploadMbps"`
}

func startHttpServer(serverPort *int) error {
	// Create a sub filesystem for the 'static' directory
	subFS, err := fs.Sub(staticFiles, "static")
	if err != nil {
		return fmt.Errorf("Failed to create sub filesystem: %v", err)
	}

	// Set up the HTTP server
	http.HandleFunc("/metrics", handleMetrics)
	http.Handle("/", http.FileServer(http.FS(subFS)))

	// Start the server with the specified port
	return http.ListenAndServe(fmt.Sprintf(":%d", *serverPort), nil)
}

func handleMetrics(w http.ResponseWriter, _ *http.Request) {
	metrics, err := getAllMetrics()
	if err != nil {
		logger.Errorf("Error loading data: %v", err)
		http.Error(w, "Error loading data", http.StatusInternalServerError)
		return
	}

	dtos := make([]MetricsDto, 0, len(metrics))
	for _, m := range metrics {
		if m.Data == nil {
			dtos = append(dtos, MetricsDto{
				EpochMillis: m.Timestamp.UnixMilli(),
				NetworkName: m.NetworkName,
				IsOnline:    m.Online,
			})
		} else {
			dtos = append(dtos, MetricsDto{
				EpochMillis:          m.Timestamp.UnixMilli(),
				NetworkName:          m.NetworkName,
				IsOnline:             m.Online,
				PingMillis:           m.Data.Latency.Milliseconds(),
				JitterMillis:         m.Data.Jitter.Milliseconds(),
				PacketLossPercentage: max(m.Data.PacketLoss.LossPercent(), 0),
				DownloadMbps:         m.Data.DLSpeed.Mbps(),
				UploadMbps:           m.Data.ULSpeed.Mbps(),
			})
		}
	}

	clientDataJson, err := json.Marshal(dtos)
	if err != nil {
		logger.Errorf("Error encoding metrics: %v", err)
		http.Error(w, "Error encoding metrics", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Connection", "close") // Close the connection after sending the response
	_, err = w.Write(clientDataJson)
	if err != nil {
		logger.Errorf("Error writing response: %v", err)
	}
}
