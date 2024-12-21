package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

type ClientData struct {
	EpochMillis          int64   `json:"EpochMillis"`
	NetworkName          string  `json:"NetworkName"`
	IsOnline             bool    `json:"IsOnline"`
	PingMillis           int64   `json:"PingMillis"`
	JitterMillis         int64   `json:"JitterMillis"`
	PacketLossPercentage float64 `json:"PacketLossPercentage"`
	DownloadMbps         float64 `json:"DownloadMbps"`
	UploadMbps           float64 `json:"UploadMbps"`
}

func startHttpServer(port *int) {
	// Set up the HTTP server
	http.HandleFunc("/metrics", handleMetrics)
	http.Handle("/", http.FileServer(http.Dir("static")))

	// Start the server with the specified port
	log.Printf("Starting server on :%d\n", *port)
	err := http.ListenAndServe(fmt.Sprintf(":%d", *port), nil)
	if err != nil {
		log.Fatalf("Error starting server: %v\n", err)
	} else {
		log.Println("Server stopped")
	}
}

func handleMetrics(w http.ResponseWriter, _ *http.Request) {
	metrics, err := getAllMetrics()
	if err != nil {
		http.Error(w, "Error loading data", http.StatusInternalServerError)
		return
	}

	clientData := make([]ClientData, 0, len(metrics))
	for _, m := range metrics {
		packetLossPercentage := m.PacketLossPercentage
		if packetLossPercentage < 0 {
			packetLossPercentage = 0
		}
		clientData = append(clientData, ClientData{
			EpochMillis:          m.Timestamp.UnixMilli(),
			NetworkName:          m.NetworkName,
			IsOnline:             m.Online,
			PingMillis:           m.PingNanos.Milliseconds(),
			JitterMillis:         m.JitterNanos.Milliseconds(),
			PacketLossPercentage: packetLossPercentage,
			DownloadMbps:         m.DownloadMbps,
			UploadMbps:           m.UploadMbps,
		})
	}

	clientDataJson, err := json.Marshal(clientData)
	if err != nil {
		http.Error(w, "Error encoding metrics", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	write, err := w.Write(clientDataJson)
	if err != nil {
		log.Printf("Error writing response: %v\n", err)
	} else {
		log.Printf("Wrote %d bytes\n", write)
	}
}
